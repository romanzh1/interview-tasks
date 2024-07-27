package main

import (
	"context"
	"fmt"
)

func sink(ctx context.Context, in <-chan int64) (chan int64, <-chan error) {
	out := make(chan int64, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		defer close(out)

		var sum int64
		for n := range in {
			select {
			case <-ctx.Done():
				return
			default:
				if n > 100 {
					errc <- fmt.Errorf("%d is bigger than 100", n)
					return
				}
				sum += n
			}
		}
		out <- sum
	}()

	return out, errc
}
