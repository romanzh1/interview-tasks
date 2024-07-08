FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o cartbuild ./cart/cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cartbuild /app/cartbuild

CMD ["/app/cartbuild"]
