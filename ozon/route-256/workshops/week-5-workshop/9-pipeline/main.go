package main

import (
	"context"
	"fmt"
	"github.com/deliveryhero/pipeline/v2"
	"strconv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		step1 = pipeline.Emit("1", "12", "23", "150")
		sum   int64
	)

	parser := pipeline.NewProcessor(func(_ context.Context, s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}, func(i string, err error) {
		fmt.Printf("can't parse value: %s, %v", i, err)
	})

	sinker := pipeline.NewProcessor(func(ctx context.Context, i int64) (int64, error) {
		if i > 100 {
			return 0, fmt.Errorf("%d is bigger than 100\n", i)
		}
		sum += i
		return sum, nil
	}, func(i int64, err error) {
		fmt.Printf("can't sum %v", err)
	})

	var (
		step2 = pipeline.Process(ctx, parser, step1)
		step3 = pipeline.Process(ctx, sinker, step2)
	)
	fmt.Println("started")

	pipeline.Drain(step3)
	fmt.Printf("result is: %d\n", sum)

}
