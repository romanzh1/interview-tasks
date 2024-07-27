package cache

import (
	"context"
	"fmt"
	logger_custom "gitlab.ozon.dev/12/week-6-workshop/pkg/logger"
	"go.opentelemetry.io/otel"
	"sync"
	"time"
)

type LeakyCache struct {
	mx      sync.Mutex
	storage map[string][]byte
}

func NewCache(size int) *LeakyCache {
	return &LeakyCache{
		storage: make(map[string][]byte, size),
	}
}

func (c *LeakyCache) Set(ctx context.Context, key string, value []byte) error {
	logger_custom.Infow(ctx, "method set is called", "key", key)

	c.mx.Lock()
	c.storage[key] = value
	c.mx.Unlock()

	buf := make([]byte, 10<<20)

	go func() {
		createdAt := time.Now()
		for {
			buf = append(buf, make([]byte, 10<<20)...)
			fmt.Println("buf len: ", len(buf))
			if time.Since(createdAt).Minutes() > 10 {
				c.mx.Lock()
				delete(c.storage, key)
				c.mx.Unlock()

				return
			}

			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}

func (c *LeakyCache) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, span := otel.GetTracerProvider().Tracer("workshop-6").Start(ctx, "cache_Get")
	defer span.End()

	c.mx.Lock()
	defer c.mx.Unlock()

	return c.storage[key], nil
}
