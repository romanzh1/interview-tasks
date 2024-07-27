package main

import (
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka"
)

type config struct {
	KafkaConfig kafka.Config
}

func newConfig() config {
	return config{
		KafkaConfig: kafka.Config{
			Brokers: []string{
				"localhost:9092",
			},
		},
	}
}
