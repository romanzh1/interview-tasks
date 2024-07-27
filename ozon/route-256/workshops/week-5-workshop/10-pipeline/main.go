package main

import (
	"context"
	"fmt"
	"github.com/deliveryhero/pipeline/v2"
	"github.com/schollz/progressbar/v3"
	"io"
	"math"
	"net/http"
	"sort"
)

func getUrlContent(link string) (result []byte, err error) {
	var resp *http.Response
	if resp, err = http.Get(link); err != nil {
		return nil, fmt.Errorf("get %s: %w", link, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s: resp code: %d", link, resp.StatusCode)
	}

	if result, err = io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("get %s: read body: %w", link, err)
	}
	return
}

type link struct {
	name     string
	link     string
	category string
	stars    int
}

var categories = []string{"Distributed Systems", "Goroutines", "Hardware"}

func main() {
	var (
		ctx         = context.Background()
		accumulator = make(map[string][]*link)
		bar         = progressbar.Default(-1)

		downloader = pipeline.NewProcessor(func(ctx context.Context, i *link) (*link, error) {
			stars, sErr := getStars(ctx, i.link)
			if sErr != nil {
				return nil, sErr
			}
			i.stars = stars
			return i, nil
		}, func(i *link, err error) {
			fmt.Printf("can't get stars count for repo: %s, %v\n", i.name, err)
		})

		progress = pipeline.NewProcessor(func(ctx context.Context, i *link) (*link, error) {
			_ = bar.Add(1)
			return i, nil
		}, nil)

		sinker = pipeline.NewProcessor(func(ctx context.Context, i *link) (*link, error) {
			accumulator[i.category] = append(accumulator[i.category], i)
			return i, nil
		}, nil)

		step1 = source(ctx, categories)
		step2 = pipeline.ProcessConcurrently(ctx, 200, downloader, step1)
		// step2 = pipeline.Process(ctx, downloader, step1)
		step3 = pipeline.Process(ctx, progress, step2)
		step4 = pipeline.Process(ctx, sinker, step3)
	)
	pipeline.Drain(step4)

	for category, categoryLinks := range accumulator {
		sort.Slice(categoryLinks, func(i, j int) bool {
			return categoryLinks[i].stars > categoryLinks[j].stars
		})
		fmt.Println()
		fmt.Println(category)
		lim := math.Min(float64(len(categoryLinks)), 20)
		for i := 0; i < int(lim); i++ {
			fmt.Printf("\t %s (%s): %d\n", categoryLinks[i].name, categoryLinks[i].link, categoryLinks[i].stars)
		}
	}

}
