package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"week-4-workshop/cart/internal/domain"
)

type Client struct {
	token    string
	basePath string
}

type GetProductRequest struct {
	Token string `json:"token,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
}

type GetProductResponse struct {
	Name  string `json:"name,omitempty"`
	Price uint32 `json:"price,omitempty"`
}

type GetProductErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

const handlerName = "get_product"

func New(basePath string, token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("product service has empty auth token")
	}

	return &Client{
		token:    token,
		basePath: basePath,
	}, nil
}

func (c Client) GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error) {
	request := GetProductRequest{
		Token: c.token,
		SKU:   sku,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request %w", err)
	}

	path, err := url.JoinPath(c.basePath, handlerName)
	if err != nil {
		return nil, fmt.Errorf("incorrect base basePath for %q: %w", handlerName, err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if httpResponse.StatusCode != http.StatusOK {
		response := &GetProductErrorResponse{}
		err = json.NewDecoder(httpResponse.Body).Decode(response)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, fmt.Errorf("HTTP request responded with: %d , message: %s", httpResponse.StatusCode, response.Message)
	}

	response := &GetProductResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &domain.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
