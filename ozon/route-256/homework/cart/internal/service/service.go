package service

import (
	"fmt"
	"sort"

	"route256/cart/internal/models"
	"route256/cart/pkg/product"
)

//go:generate minimock -i cartRepository -o ./mocks/ -s _mock.go
//go:generate minimock -i productClient -o ./mocks/ -s _mock.go
type cartRepository interface {
	AddItemToUserCart(userID, skuID int64, count uint16) error
	DeleteItemFromUserCart(userID, skuID int64) error
	ClearUserCart(userID int64) error
	ListUserCart(userID int64) ([]models.CartItem, error)
}

type productClient interface {
	GetProduct(skuID uint32) (*product.Product, error)
}

type Service struct {
	repo          cartRepository
	productClient productClient
}

func NewService(repo cartRepository, productClient productClient) *Service {
	return &Service{repo: repo, productClient: productClient}
}

func (s *Service) AddItemToUserCart(userID, skuID int64, count uint16) error {
	checkProduct, err := s.productClient.GetProduct(uint32(skuID))
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if checkProduct == nil {
		return fmt.Errorf("product not found")
	}

	return s.repo.AddItemToUserCart(userID, skuID, count)
}

func (s *Service) DeleteItemFromUserCart(userID, skuID int64) error {
	return s.repo.DeleteItemFromUserCart(userID, skuID)
}

func (s *Service) ClearUserCart(userID int64) error {
	return s.repo.ClearUserCart(userID)
}

func (s *Service) ListUserCart(userID int64) ([]models.CartItem, uint32, error) {
	items, err := s.repo.ListUserCart(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cart: %w", err)
	}

	var cartItems []models.CartItem
	var totalPrice uint32

	for _, item := range items {
		p, err := s.productClient.GetProduct(uint32(item.SkuID))
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
