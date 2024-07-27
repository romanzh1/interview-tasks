package consumer_group

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type consumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (c *consumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[consumer-group] run")

		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
				log.Printf("Error from consume: %v\n", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Printf("[consumer-group]: ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

}

func NewConsumerGroup(brokers []string, groupID string, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...Option) (*consumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	config.Consumer.Group.ResetInvalidOffsets = true
	// Сердцебиение консьюмера
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	// Таймаут сессии
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	// Таймаут ребалансировки
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	//
	config.Consumer.Return.Errors = true

	//config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Применяем свои конфигурации
	for _, opt := range opts {
		opt.Apply(config)
	}

	/*
	  Setup a new Sarama consumer group
	*/
	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}
