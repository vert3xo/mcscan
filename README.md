# Minecraft Scanner
Program that scans the entire internet for servers with default Minecraft port open (25565).

It saves the results in MongoDB database. Mongo Express frontend is also available.

# Usage
## Docker
Simplest way to run this project is with Docker:
`docker compose up`

## Compiling and running manualy
This project contains hardcoded IP addresses which Go allows us to change during the compilation.
If you are running both Redis and Mongo on localhost and default port you can run `make all`
### Environment variables
Except the addresses all other values can be configured in `.env`. To get started just rename `.env.example` to `.env` and change the values.
