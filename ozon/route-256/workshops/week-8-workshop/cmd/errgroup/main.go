package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

// код взят из примера https://pkg.go.dev/golang.org/x/sync/errgroup
func main() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		// ! конкурентный доступ к элементу массива без мьютекса
		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for i, search := range searches {
			// c go 1.22 это можно уже не делать
			i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := search(ctx, query)
				if err == nil {
					// ! конкурентный доступ к элементу массива без мьютекса
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
