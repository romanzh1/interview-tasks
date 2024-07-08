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

.PHONY: proto
proto:
	protoc --proto_path=proto \
           --proto_path=. \
           --go_out=proto --go_opt=paths=source_relative \
           --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
           --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative,logtostderr=true \
           proto/*.proto

.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

.PHONY: run
