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
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

gocheck:
	gocyclo -top 10 -ignore "mocks" .
	gocognit -ignore "mock|_test" .

.PHONY: proto
proto:
	protoc --proto_path=proto \
           --proto_path=. \
           --go_out=proto --go_opt=paths=source_relative \
           --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
           --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative,logtostderr=true \
           proto/*.proto

.PHONY: run
