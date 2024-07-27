package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func makeGetRequest(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("url: %s | create request %w", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("url: %s | do request: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url: %s | http status code: %s", url, resp.Status)
	}
	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() (err error) {
		defer func() {
			fmt.Println(err)
		}()
		return makeGetRequest(ctx, "https://profit.ozon.ru")
	})
	eg.Go(func() (err error) {
		defer func() {
			fmt.Println(err)
		}()
		return makeGetRequest(ctx, "https://ytrtrtr454a.ru")
	})
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("done")

}
