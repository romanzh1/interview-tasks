include .env
export

run:
	go run ./cmd

lint:
	gofumpt -l -w .
	find . -name '*.go' -exec gci write --skip-generated -s standard -s default {} \;
	find . -name '*.go' | xargs goimports -w
	golangci-lint run ./...

test: lint
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

gocheck:
	gocyclo -top 10 -ignore "mocks" .
	gocognit -ignore "mock|_test" .

.PHONY: run
