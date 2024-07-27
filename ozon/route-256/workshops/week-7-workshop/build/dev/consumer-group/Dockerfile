# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY ./ ./

RUN pwd && ls -alg && CGO_ENABLED=0 GOOS=linux go build ./cmd/consumer_group

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/consumer_group /consumer_group

# EXPOSE 8080

#USER nonroot:nonroot
ENTRYPOINT ["/consumer_group", "-bootstrap-server=kafka0:29092"]
