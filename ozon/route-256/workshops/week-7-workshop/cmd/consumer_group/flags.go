package main

import "flag"

type flags struct {
	topic             string
	bootstrapServer   string
	consumerGroupName string
}

func init() {
	flag.StringVar(&cliFlags.topic, "topic", "route256-example", "topic to produce")
	flag.StringVar(&cliFlags.bootstrapServer, "bootstrap-server", "localhost:9092", "kafka broker host and port")
	flag.StringVar(&cliFlags.consumerGroupName, "cg-name", "route256-consumer-group", "topic to produce")

	flag.Parse()
}
