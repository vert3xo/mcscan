# Minecraft Scanner
Program that scans the entire internet for servers with default Minecraft port open (25565).

It saves the results in MongoDB database. Mongo Express frontend is also available.

# Known Issues
The masscan wrapper library this project uses currently doesn't support exclusion of IP addresses and thus doesn't support scanning of range `0.0.0.0/0`

A pull request to fix this is currently pending. I will update this project as soon as this problem is fixed or create a custom wrapper myself.

# Usage
## Environment variables
Except the addresses all other values can be configured in `.env`.

To get started just rename `.env.example` to `.env` and change the values.

## Docker
Simplest way to run this project is with Docker:

`docker compose up`

## Compiling and running manualy
### The easy way
This project contains hardcoded IP addresses which Go allows us to change during the compilation.

If you are running both Redis and Mongo on localhost and default port you can run `make all` to compile everything.

To compile only the scanner or only the worker you can run `make scanner` or `make server` respectively.

### More complicated way
If you need to modify any addresses you can use the following commands for compilation:

#### Worker
`go build -o worker -ldflags='-X main.redisAddr=<redis_address>:<redis_port> -X main.mongoAddr=<mongo_address>:<mongo_port>' worker/worker.go`

#### Scanner
`go build -o scanner -ldflags='-X main.redisAddr=<redis_address>:<redis_port>' main.go`