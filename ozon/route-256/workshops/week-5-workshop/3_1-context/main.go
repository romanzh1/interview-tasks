package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	start := time.Now()
	val, err := operate(ctx)
	fmt.Printf("got %d, %v after %s\n", val, err, time.Since(start))
}

func operate(ctx context.Context) (int64, error) {
	var out int64

	for i := range 1_000_000_000 {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			out += int64(i)
		}
	}

	return out, nil
}
