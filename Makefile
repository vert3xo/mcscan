all:
	make client
	make server 

client:
	go build -v -o client main.go

server:
	go build -v -o worker/worker worker/worker.go

clean:
	rm -f client worker/worker
