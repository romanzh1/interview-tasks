package main

import "flag"

type flags struct {
	topic string
}

func init() {
	flag.StringVar(&cliFlags.topic, "topic", "route256-example", "topic to produce")

	flag.Parse()
}
