package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alteamc/minequery/ping"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client    *mongo.Client
	Ctx       context.Context
	CtxCancel context.CancelFunc
	Error     error
}

var database Database

const (
	TypeGetPingResponse = "ping:getresponse"
)

type ServerPayload struct {
	Ip   string
	Port int
}

type PlayersList struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

type PingResponse struct {
	Host          string        `json:"host"`
	Port          int           `json:"port"`
	Version       string        `json:"version"`
	MaxPlayers    int           `json:"maxPlayers"`
	OnlinePlayers int           `json:"onlinePlayers"`
	PlayersList   []PlayersList `json:"playersList"`
	Description   interface{}   `json:"description"`
	Favicon       string        `json:"favicon"`
}

func ConnectMongo(mongoUsername, mongoPassword, mongoSocket, authDB string) Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoUsername, mongoPassword, mongoSocket, authDB)))
	if err != nil {
		database = Database{Client: nil, Ctx: nil, CtxCancel: nil, Error: err}
		return database
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*20)
	err = client.Connect(ctx)
	if err != nil {
		database = Database{Client: nil, Ctx: ctx, CtxCancel: ctxCancel, Error: err}
		return database
	}

	database = Database{Client: client, Ctx: ctx, CtxCancel: ctxCancel, Error: nil}
	return database
}

func intoPingResponse(host string, port int, response ping.Response) *PingResponse {
	playersList := []PlayersList{}
	for i := 0; i < len(response.Players.Sample); i++ {
		playersList = append(playersList, PlayersList{
			UUID:     response.Players.Sample[i].ID,
			Username: response.Players.Sample[i].Name,
		})
	}

	return &PingResponse{
		Host:          host,
		Port:          port,
		Version:       response.Version.Name,
		MaxPlayers:    response.Players.Max,
		OnlinePlayers: response.Players.Online,
		PlayersList:   playersList,
		Description:   response.Description,
		Favicon:       response.Favicon,
	}
}

func NewGetPingResponseTask(ip string, port int) (*asynq.Task, error) {
	payload, err := json.Marshal(ServerPayload{ip, port})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeGetPingResponse, payload), nil
}

func HandleGetPingResponseTask(ctx context.Context, t *asynq.Task) error {
	var p ServerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal the payload: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Pinging %s:%d", p.Ip, p.Port)
	res, err := ping.PingWithTimeout(p.Ip, uint16(p.Port), time.Second*10)
	if err == nil {
		properResponse := intoPingResponse(p.Ip, p.Port, *res)
		_, err = database.Client.Database("mcscan").Collection("servers").InsertOne(database.Ctx, properResponse)
		if err != nil {
			log.Printf("failed to insert data: %v", err)
		}
	} else {
		return fmt.Errorf("server does not look like a minecraft server")
	}

	return nil
}
