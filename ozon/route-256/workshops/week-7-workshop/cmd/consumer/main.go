package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka/consumer"
)

var cliFlags = flags{}

func main() {
	var (
		wg   = &sync.WaitGroup{}
		conf = newConfig()
		ctx  = runSignalHandler(context.Background(), wg)
	)

	cons, err := consumer.NewConsumer(conf.KafkaConfig,
		consumer.WithReturnErrorsEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer cons.Close() // Не забываем освобождать ресурсы :)

	err = cons.ConsumeTopic(ctx, cliFlags.topic, func(msg *sarama.ConsumerMessage) {
		// Your logic here

		data, _ := json.Marshal(convertMsg(msg))
		fmt.Printf("Read Topic: %s\n", data)
	}, wg)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	//cons.Close() // call after all consumers already closed
}

func runSignalHandler(ctx context.Context, wg *sync.WaitGroup) context.Context {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer signal.Stop(sigterm)
		defer wg.Done()
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					fmt.Printf("[signal] signal chan closed: %s\n", sig.String())
					return
				}

				fmt.Printf("[signal] signal recv: %s\n", sig.String())
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					fmt.Println("[signal] context closed")
					return
				}

				fmt.Printf("[signal] ctx done: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

	return sigCtx
}

type Msg struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func convertMsg(in *sarama.ConsumerMessage) Msg {
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       string(in.Key),
		Payload:   string(in.Value),
	}
}
