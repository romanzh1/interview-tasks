package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()

	// write value
	ctx = context.WithValue(ctx, "key", "value")
	// ctx = context.WithValue(ctx, "key", "value1")

	// ВНИМАНИЕ!
	// Читающий эти строки - не клади транзакции в контекст
	// Смотрящий в этот код - не используй контекст для неявной передачи аргументов функции
	v := ctx.Value("key")
	fmt.Println(v.(string))

	// Контракт value оставляет желать лучшего
	// fmt.Println(v.(int))

	ctx, cancel := context.WithCancel(ctx)

	from := time.Now()

	go func() {
		for {
			fmt.Println("cancel")
			cancel()
		}
	}()

	<-ctx.Done()
	fmt.Println("close app after", time.Since(from))
}
