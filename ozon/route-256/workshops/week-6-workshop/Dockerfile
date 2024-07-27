FROM gitlab-registry.ozon.dev/go/classroom-12/students/homework/golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/server

FROM gitlab-registry.ozon.dev/go/classroom-12/students/homework/alpine

COPY --from=builder /app/bin/server /server

CMD [ "/server" ]
