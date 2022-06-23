<h1 align="center">Minecraft Scanner</h1>
<h3 align="center">A tool that scans the Internet for open Minecraft servers</h3>

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

`go build -o worker worker/worker.go`

#### Scanner

`go build -o scanner main.go`
