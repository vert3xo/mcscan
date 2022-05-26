package main

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/vert3xo/mcscan/tasks"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	redisAddr := os.Getenv("REDIS_SOCKET")
	mongoAddr := os.Getenv("MONGO_SOCKET")

	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	database := tasks.ConnectMongo(mongoUsername, mongoPassword, mongoAddr, "admin")
	if database.Error != nil {
		log.Fatalf("failed to connect to the database: %v", database.Error)
	}
	defer database.Client.Disconnect(database.Ctx)

	server := asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr}, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeGetPingResponse, tasks.HandleGetPingResponseTask)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
