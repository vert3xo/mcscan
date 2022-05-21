all:
	make scanner
	make work

scanner:
	go build -v -o client main.go

work:
	go build -v -o worker/worker worker/worker.go

clean:
	rm -f client worker/worker
