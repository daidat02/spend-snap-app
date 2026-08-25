.PHONY: run build tidy test lint migrate-up migrate-down

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

lint:
	go vet ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1
