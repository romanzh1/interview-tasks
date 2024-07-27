package main

import (
	"context"
	"strconv"
)

func parse(ctx context.Context, input <-chan string) (<-chan int64, <-chan error) {
	var (
		out  = make(chan int64)
		errc = make(chan error, 1)
	)

	go func() {
		defer close(out)
		defer close(errc)

		for line := range input {
			n, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				errc <- err
				return
			}

			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc
}
