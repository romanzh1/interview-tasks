package producer

import (
	"time"

	"github.com/IBM/sarama"
)

// Option is a configuration callback
type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

// WithProducerPartitioner ...
func WithProducerPartitioner(pfn sarama.PartitionerConstructor) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Partitioner = pfn
		return nil
	})
}

// WithProducerPartitioner ...
func WithRequiredAcks(acks sarama.RequiredAcks) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.RequiredAcks = acks
		return nil
	})
}

// WithIdempotent ...
func WithIdempotent() Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Idempotent = true
		return nil
	})
}

// WithMaxRetries ...
func WithMaxRetries(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Max = n
		return nil
	})
}

// WithRetryBackoff ...
func WithRetryBackoff(d time.Duration) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Backoff = d
		return nil
	})
}

// WithMaxOpenRequests ...
func WithMaxOpenRequests(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Net.MaxOpenRequests = n
		return nil
	})
}

// WithProducerFlushMessages ...
func WithProducerFlushMessages(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Flush.Messages = n
		return nil
	})
}

// WithProducerFlushFrequency ...
func WithProducerFlushFrequency(d time.Duration) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Flush.Frequency = d
		return nil
	})
}
