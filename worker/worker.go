package main

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/vert3xo/mcscan/tasks"
)

var redisAddr = "127.0.0.1:6379"
var mongoAddr = "127.0.0.1:27017"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	database := tasks.ConnectMongo(mongoUsername, mongoPassword, mongoAddr, "admin")
	if database.Error != nil {
		log.Fatalf("failed to connect to the database: %v", database.Error)
	}
	defer database.Client.Disconnect(database.Ctx)
	defer database.CtxCancel()

	server := asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr}, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeGetPingResponse, tasks.HandleGetPingResponseTask)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}