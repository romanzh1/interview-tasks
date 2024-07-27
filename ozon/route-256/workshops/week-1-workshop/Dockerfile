FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/server.go

FROM scratch
COPY --from=builder server /bin/server

ENTRYPOINT ["/bin/server"]