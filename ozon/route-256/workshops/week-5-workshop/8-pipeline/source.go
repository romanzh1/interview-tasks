package main

import (
	"context"
	"errors"
	"fmt"
)

func source(ctx context.Context, lines ...string) (<-chan string, <-chan error, error) {
	if len(lines) == 0 {
		return nil, nil, errors.New("empty array")
	}

	var (
		out  = make(chan string)
		errc = make(chan error, 1)
	)

	go func() {
		defer close(out)
		defer close(errc)

		for idx, line := range lines {
			if len(line) == 0 {
				errc <- fmt.Errorf("line %d is empty", idx)
				return
			}

			select {
			case out <- line:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc, nil
}
