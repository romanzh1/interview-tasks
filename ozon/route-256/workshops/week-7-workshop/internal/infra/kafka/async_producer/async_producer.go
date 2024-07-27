package async_producer

import (
	"fmt"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/go/classroom-12/students/week-7-workshop/internal/infra/kafka"
)

func NewAsyncProducer(conf kafka.Config, opts ...Option) (sarama.AsyncProducer, error) {
	config := PrepareConfig(opts...)

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return asyncProducer, nil
}
