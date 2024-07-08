FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o lomsbuild ./loms/cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lomsbuild /app/lomsbuild

CMD ["/app/lomsbuild"]
