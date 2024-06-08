include .env
export

run:
	go run ./cmd

lint:
	golangci-lint run ./...

test: lint
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: run
