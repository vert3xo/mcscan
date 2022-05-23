package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/vert3xo/mcscan/tasks"
	"github.com/zan8in/masscan"
)

var redisAddr = "127.0.0.1:6379"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	rate, _ := strconv.Atoi(os.Getenv("SCAN_RATE"))

	scanner, err := masscan.NewScanner(
		masscan.SetParamRate(rate),
		masscan.SetParamTargets(os.Getenv("SCAN_RANGE")),
		masscan.SetParamPorts("25565"),
		masscan.SetParamWait(5),
		masscan.EnableDebug(),
	)

	if err != nil {
		log.Fatalf("failed to create a scanner: %v", err)
	}

	if err = scanner.RunAsync(); err != nil {
		log.Fatalf("failed to start the scanner: %v", err)
	}

	stdout := scanner.GetStdout()

	go func() {
		for stdout.Scan() {
			res := masscan.ParseResult(stdout.Bytes())

			port, _ := strconv.Atoi(res.Port)

			task, err := tasks.NewGetPingResponseTask(res.IP, port)
			if err != nil {
				log.Fatalf("could not create task: %v", err)
			}

			info, err := client.Enqueue(task)
			if err != nil {
				log.Fatalf("could not enqueue task: %v", err)
			}

			log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
		}
	}()

	if err = scanner.Wait(); err != nil {
		log.Fatalf("failed to wait for the scan to finish: %v", err)
	}
}
