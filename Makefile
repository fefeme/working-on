build:
	go build -o bin/wo -ldflags="-s -w" main.go

run:
	go run main.go
