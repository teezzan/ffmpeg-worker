serve:
	go run cmd/server/main.go

build:
	go build -o bin/main cmd/worker/main.go

run:
	go run cmd/worker/main.go