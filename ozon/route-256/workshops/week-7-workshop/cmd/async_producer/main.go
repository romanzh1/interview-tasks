package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/domain/order"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka/async_producer"
)

var cliFlags = flags{}

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	conf := newConfig()
	prod, err := async_producer.NewAsyncProducer(conf.KafkaConfig,
		//async_producer.WithIdempotent(),
		async_producer.WithRequiredAcks(sarama.WaitForAll),
		//async_producer.WithMaxOpenRequests(1),
		async_producer.WithMaxRetries(5),
		async_producer.WithRetryBackoff(10*time.Millisecond),
		//async_producer.WithProducerPartitioner(sarama.NewManualPartitioner),
		//async_producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		//async_producer.WithProducerPartitioner(sarama.NewRandomPartitioner),
		async_producer.WithProducerFlushMessages(8),
		async_producer.WithProducerFlushFrequency(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	runTicker(ctx, &wg)
	runKafkaHandlers(ctx, prod, &wg)

	var event order.Event
	for i := 0; i < cliFlags.repeatCnt; i++ {
		factory := order.NewDefaultFactory(cliFlags.startID)
		for i := 0; i < cliFlags.count; i++ {
			event = factory.Create(order.EventOrderCreated)

			bytes, err := json.Marshal(event)
			if err != nil {
				log.Fatal(err)
			}

			msg := &sarama.ProducerMessage{
				Topic: cliFlags.topic,
				Key:   sarama.StringEncoder(strconv.FormatInt(event.ID, 10)),
				Value: sarama.ByteEncoder(bytes),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("app-name"),
						Value: []byte("route256-async-prod"),
					},
				},
				Timestamp: time.Now(),
			}

			prod.Input() <- msg
			log.Printf("[produce] msg sent: %d\n", event.ID)
			time.Sleep(200 * time.Millisecond)
		}
	}

	prod.AsyncClose()
	<-prod.Successes()
	<-prod.Errors()

	cancel()
	wg.Wait()
}

func runTicker(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fmt.Println("[ticker] start")
		defer wg.Done()
		tickN := 0

		ticker := time.NewTicker(time.Millisecond * 100)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[ticker] terminate")
				return
			case <-ticker.C:
				tickN++
				log.Printf("[ticker] %d\n", tickN*100)
			}
		}
	}()
}

// !!!ВНИМАНИЕ!!!
// ОБЯЗАТЕЛЬНОЕ чтение канала успешных событий при c.Producer.Return.Successes = true
func runKafkaHandlers(ctx context.Context, prod sarama.AsyncProducer, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fmt.Println("[sent success/error] start")
		defer wg.Done()

		successCh := prod.Successes()
		errCh := prod.Errors()

		for {
			select {
			case <-ctx.Done():
				log.Println("[sent success/error] terminate")
				return
			case msg := <-successCh:
				if msg == nil {
					log.Println("[sent success] chan closed")
					return
				}
				log.Printf("[sent success] key: %q, partition: %d, offset: %d\n", msg.Key, msg.Partition, msg.Offset)
			case msgErr := <-errCh:
				if msgErr == nil {
					log.Println("[sent error] chan closed")
					return
				}
				log.Printf("[sent error] err %s, topic: %q, offset: %d\n", msgErr.Err, msgErr.Msg.Topic, msgErr.Msg.Offset)
			}
		}
	}()
}
