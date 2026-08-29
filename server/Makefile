-include .env
export
DATABASE_URL ?= postgres://postgres:secretpassword@localhost:5432/spendsnap_db?sslmode=disable

.PHONY: run build tidy test lint migrate-up migrate-down migrate-create

dev:
	air

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
	
# 	ví dụ: make migrate-create name=name_table
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)