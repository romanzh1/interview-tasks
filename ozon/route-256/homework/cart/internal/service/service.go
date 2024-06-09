package service

import (
	"context"
	"fmt"
	"sort"

	"route256/cart/internal/models"
	"route256/cart/pkg/product"
)

//go:generate minimock -i cartRepository -o ./mocks/ -s _mock.go
//go:generate minimock -i productClient -o ./mocks/ -s _mock.go
type cartRepository interface {
	AddItemToUserCart(ctx context.Context, cart models.CartRequest) error
	DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error
	ClearUserCart(ctx context.Context, userID int64) error
	ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, error)
}

type productClient interface {
	GetProduct(ctx context.Context, skuID uint32) (*product.Product, error)
}

type Service struct {
	repo          cartRepository
	productClient productClient
}

func NewService(repo cartRepository, productClient productClient) *Service {
	return &Service{repo: repo, productClient: productClient}
}

func (s *Service) AddItemToUserCart(ctx context.Context, cart models.CartRequest) error {
	checkProduct, err := s.productClient.GetProduct(ctx, uint32(cart.SkuID))
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if checkProduct == nil {
		return fmt.Errorf("product not found")
	}

	return s.repo.AddItemToUserCart(ctx, cart)
}

func (s *Service) DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error {
	return s.repo.DeleteItemFromUserCart(ctx, userID, skuID)
}

func (s *Service) ClearUserCart(ctx context.Context, userID int64) error {
	return s.repo.ClearUserCart(ctx, userID)
}

func (s *Service) ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, uint32, error) {
	items, err := s.repo.ListUserCart(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cart: %w", err)
	}

	cartItems := make([]models.CartItem, 0, len(items))
	totalPrice := uint32(0)

	for _, item := range items {
		p, err := s.productClient.GetProduct(ctx, uint32(item.SkuID))
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get p: %w", err)
		}

		if p == nil {
			continue
		}

		cartItems = append(cartItems, models.CartItem{
			SkuID: item.SkuID,
			Name:  p.Name,
			Count: item.Count,
			Price: p.Price,
		})

		totalPrice += p.Price * uint32(item.Count)
	}

	sort.Slice(cartItems, func(i, j int) bool {
		return cartItems[i].SkuID < cartItems[j].SkuID
	})

	return cartItems, totalPrice, nil
}
