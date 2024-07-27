package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/domain/order"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka/producer"
)

func main() {
	conf := newConfig(cliFlags)
	fmt.Printf("%+v\n", conf)
	prod, err := producer.NewSyncProducer(conf.kafka,
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		//producer.WithProducerPartitioner(sarama.NewManualPartitioner),
		//producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		//producer.WithProducerPartitioner(sarama.NewRandomPartitioner),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer prod.Close() // Не забываем освобождать ресурсы

	var event order.Event
	for i := 0; i < conf.app.repeatCnt; i++ {
		factory := order.NewDefaultFactory(conf.app.startID)
		for i := 0; i < conf.app.count; i++ {
			event = factory.Create(order.EventOrderCreated)

			bytes, err := json.Marshal(event)
			if err != nil {
				log.Fatal(err)
			}

			msg := &sarama.ProducerMessage{
				Topic: conf.producer.topic,
				Key:   sarama.StringEncoder(strconv.FormatInt(event.ID, 10)),
				Value: sarama.ByteEncoder(bytes),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("app-name"),
						Value: []byte("route256-sync-prod"),
					},
				},
				//Partition: 1,
				Timestamp: time.Now(),
			}

			partition, offset, err := prod.SendMessage(msg)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("key: %d, partition: %d, offset: %d", event.ID, partition, offset)

			time.Sleep(conf.app.interval)
		}
	}
}
