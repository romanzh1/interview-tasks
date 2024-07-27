package add

import (
	"context"
	"errors"
	"fmt"

	"week-2-workshop/cart/internal/domain"
)

type (
	productService interface {
		GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error)
	}
	repository interface {
		AddItem(ctx context.Context, userID int64, item domain.Item) error
	}

	Handler struct {
		productService productService
		repo           repository
	}
)

var (
	ErrInvalidSKU = errors.New("invalid sku")
)

func New(repo repository, productService productService) *Handler {
	return &Handler{
		repo:           repo,
		productService: productService,
	}
}

func (h *Handler) AddItem(ctx context.Context, userID int64, item domain.Item) error {
	products, err := h.productService.GetProductInfo(ctx, item.SKU)
	if err != nil {
		return fmt.Errorf("ProductService.GetProductInfo failed: %w", err)
	}

	if products == nil {
		return fmt.Errorf("ProductService.GetProductInfo return no product with given SKU=%d: %w", item.SKU, ErrInvalidSKU)
	}

	err = h.repo.AddItem(ctx, userID, item)
	if err != nil {
		return fmt.Errorf("repo.AddItem failed: %w", err)
	}

	return nil
}

func (h *Handler) AddItemV2(ctx context.Context, userID int64, item1, item2 domain.Item) error {
	products, err := h.productService.GetProductInfo(ctx, item1.SKU)
	if err != nil {
		return fmt.Errorf("ProductService.GetProductInfo failed: %w", err)
	}

	if products == nil {
		return fmt.Errorf("ProductService.GetProductInfo return no product with given SKU=%d: %w", item1.SKU, ErrInvalidSKU)
	}

	products, err = h.productService.GetProductInfo(ctx, item2.SKU)
	if err != nil {
		return fmt.Errorf("ProductService.GetProductInfo failed: %w", err)
	}

	if products == nil {
		return fmt.Errorf("ProductService.GetProductInfo return no product with given SKU=%d: %w", item2.SKU, ErrInvalidSKU)
	}

	err = h.repo.AddItem(ctx, userID, item1)
	if err != nil {
		return fmt.Errorf("repo.AddItem failed: %w", err)
	}

	return nil
}
