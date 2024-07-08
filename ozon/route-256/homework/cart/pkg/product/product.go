package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"route256/cart/cmd/config"
)

type Client struct {
	baseURL     string
	token       string
	client      *http.Client
	rateLimiter *time.Ticker
}

func NewClient(cfg config.ProductClient) (*Client, error) {
	productClientTimeout, err := time.ParseDuration(cfg.Timeout + "s")
	if err != nil {
		return nil, fmt.Errorf("failed to parse product client timeout: %w", err)
	}

	return &Client{
		baseURL:     cfg.Host,
		token:       cfg.Token,
		client:      &http.Client{Timeout: productClientTimeout},
		rateLimiter: time.NewTicker(time.Second / time.Duration(cfg.Rps)),
	}, nil
}

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Product
	Message string `json:"message"`
}

type Request struct {
	Token string `json:"token"`
	SkuID uint32 `json:"sku"`
}

func (c *Client) GetProduct(ctx context.Context, skuID uint32) (*Product, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.rateLimiter.C:
	}

	url := fmt.Sprintf("%s/get_product", c.baseURL)

	requestBody, err := json.Marshal(Request{
		Token: c.token,
		SkuID: skuID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get product: %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Message == "sku not found" {
		return nil, fmt.Errorf("product not found")
	}

	return &response.Product, nil
}
