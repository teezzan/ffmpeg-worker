build:
	go build -o bin/main cmd/worker/main.go

run:
	go run cmd/worker/main.go

publish:
	go run cmd/publisher/main.go