package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka"
)

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(kafkaConfig kafka.Config, opts ...Option) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Применяем свои конфигурации
	for _, opt := range opts {
		_ = opt.Apply(config)
	}

	consumer, err := sarama.NewConsumer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, err
	}

	/*
		consumer.Topics() - список топиков
		consumer.Partitions("test_topic") - партиции топика
		consumer.ConsumePartition("test_topic", 1, 12) - чтение конкретного топика с 12 сдвига в первой партиции
		consumer.Pause() - останавливаем чтение определенных топиков
		consumer.Resume() - восстанавливаем чтение определенных топиков
		consumer.PauseAll() - останавливаем чтение всех топиков
		consumer.ResumeAll() - восстанавливаем чтение всех топиков
	*/

	return &Consumer{
		consumer: consumer,
	}, err
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) ConsumeTopic(ctx context.Context, topic string, handler func(*sarama.ConsumerMessage), wg *sync.WaitGroup) error {
	// получаем все партиции топика
	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	/*
	   sarama.OffsetOldest - перечитываем каждый раз все
	   sarama.OffsetNewest - перечитываем только новые

	   Можем задавать отдельно на каждую партицию
	   Также можем сходить в отдельное хранилище и взять оттуда сохраненный offset
	*/
	var initialOffset = sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("consumer for topic=%s, partition=%d terminated\n", topic, partition)
					return
				case msg, ok := <-pc.Messages():
					if !ok {
						fmt.Println("consumer mag channel closed")
						return
					}
					handler(msg)
				}
			}
		}(pc, partition)
	}

	return nil
}
