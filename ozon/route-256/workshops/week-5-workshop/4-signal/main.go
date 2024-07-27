package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	fmt.Println("start")
	<-ctx.Done()

	// server.Close()
	// закрываем все коннекты (например, к БД)
	// коммитим офсет кафки и тд
	fmt.Println("close")
}
