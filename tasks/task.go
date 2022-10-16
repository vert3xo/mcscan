package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alteamc/minequery/ping"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vert3xo/mcscan/types"
	"github.com/vert3xo/mcscan/utils"
)

type Database struct {
	Client *mongo.Client
	Ctx    context.Context
	Error  error
}

var database Database

const (
	TypeGetPingResponse = "ping:getresponse"
)

func ConnectMongo(mongoUsername, mongoPassword, mongoSocket, authDB string) Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoUsername, mongoPassword, mongoSocket, authDB)))
	if err != nil {
		database = Database{Client: nil, Ctx: nil, Error: err}
		return database
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		database = Database{Client: nil, Ctx: ctx, Error: err}
		return database
	}

	database = Database{Client: client, Ctx: ctx, Error: nil}
	return database
}

func intoPingResponse(host string, port int, response ping.Response) *types.PingResponse {
	playersList := []types.Player{}
	for i := 0; i < len(response.Players.Sample); i++ {
		playersList = append(playersList, types.Player{
			UUID:     response.Players.Sample[i].ID,
			Username: response.Players.Sample[i].Name,
		})
	}

	return &types.PingResponse{
		Host:          	host,
		Port:          	port,
		Version:       	response.Version.Name,
		MaxPlayers:    	response.Players.Max,
		OnlinePlayers: 	response.Players.Online,
		PlayersList:	playersList,
		Description: 	response.Description,
		Favicon:       	response.Favicon,
	}
}

func NewGetPingResponseTask(ip string, port int) (*asynq.Task, error) {
	payload, err := json.Marshal(types.ServerPayload{Ip: ip, Port: port})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeGetPingResponse, payload), nil
}

func HandleGetPingResponseTask(ctx context.Context, t *asynq.Task) error {
	var p types.ServerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal the payload: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Pinging %s:%d", p.Ip, p.Port)
	res, err := ping.PingWithTimeout(p.Ip, uint16(p.Port), time.Second*10)
	if err == nil {
		properResponse := intoPingResponse(p.Ip, p.Port, *res)
		collection := database.Client.Database("mcscan").Collection("servers")
		var server bson.M
		err := collection.FindOne(database.Ctx, bson.D{primitive.E{Key: "host", Value: p.Ip}}).Decode(&server)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				_, err = collection.InsertOne(database.Ctx, properResponse)
				if err != nil {
					log.Printf("failed to insert data: %v", err)
				}
			}
		} else {
			currentPlayersList := []types.Player{}
			for _, player := range server["playerslist"].(primitive.A) {
				playerMap := player.(primitive.M)
				currentPlayersList = append(currentPlayersList, types.Player{Username: playerMap["username"].(string), UUID: playerMap["uuid"].(string)})
			}
			fmt.Println(currentPlayersList)
			fmt.Println(properResponse.PlayersList)
			properResponse.PlayersList = utils.RemoveDuplicatesFromPlayersList(append(properResponse.PlayersList, currentPlayersList...))
			update := bson.D{primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "version", Value: properResponse.Version},
				primitive.E{Key: "maxplayers", Value: properResponse.MaxPlayers},
				primitive.E{Key: "onlineplayers", Value: properResponse.OnlinePlayers},
				primitive.E{Key: "playerslist", Value: properResponse.PlayersList},
				primitive.E{Key: "description", Value: properResponse.Description},
				primitive.E{Key: "favicon", Value: properResponse.Favicon},
			}}}
			_, err := collection.UpdateOne(database.Ctx, server, update)
			if err != nil {
				return fmt.Errorf("failed to update the database: %v", err)
			}
		}
	} else {
		return fmt.Errorf("server does not look like a minecraft server")
	}

	return nil
}
