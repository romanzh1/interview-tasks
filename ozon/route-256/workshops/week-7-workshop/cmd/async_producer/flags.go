package main

import "flag"

type flags struct {
	repeatCnt int
	startID   int
	count     int
	topic     string
}

func init() {
	flag.IntVar(&cliFlags.repeatCnt, "repeat-count", 3, "count times all messages sent")
	flag.IntVar(&cliFlags.startID, "start-id", 1, "start order-id of all messages")
	flag.IntVar(&cliFlags.count, "count", 10, "count of orders to emit events")
	flag.StringVar(&cliFlags.topic, "topic", "route256-example", "topic to produce")

	flag.Parse()
}
