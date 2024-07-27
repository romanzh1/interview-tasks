package main

import (
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka"
	"time"
)

type appConfig struct {
	repeatCnt int
	startID   int
	count     int
	interval  time.Duration
}

type producerConfig struct {
	topic string
}

type config struct {
	app      appConfig
	kafka    kafka.Config
	producer producerConfig
}

func newConfig(f flags) config {
	return config{
		app: appConfig{
			repeatCnt: f.repeatCnt,
			startID:   f.startID,
			count:     f.count,
			interval:  f.interval,
		},
		kafka: kafka.Config{
			Brokers: []string{
				"localhost:9092",
			},
		},
		producer: producerConfig{
			topic: f.topic,
		},
	}
}
