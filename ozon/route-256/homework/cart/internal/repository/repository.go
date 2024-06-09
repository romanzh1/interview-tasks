package repository

import (
	"context"
	"sync"

	"route256/cart/internal/models"
)

type Repository struct {
	carts map[int64]map[int64]uint16 // userID -> skuID -> count
	mu    sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{carts: make(map[int64]map[int64]uint16)}
}

func (r *Repository) AddItemToUserCart(ctx context.Context, cart models.CartRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.carts[cart.UserID] == nil {
		r.carts[cart.UserID] = make(map[int64]uint16)
	}

	r.carts[cart.UserID][cart.SkuID] += cart.Count

	return nil
}

func (r *Repository) DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts[userID], skuID)

	return nil
}

func (r *Repository) ClearUserCart(ctx context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts, userID)

	return nil
}

func (r *Repository) ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userCart, exists := r.carts[userID]
	if !exists || len(userCart) == 0 {
		return nil, models.ErrCartIsEmpty
	}

	items := make([]models.CartItem, 0, len(userCart))
	for skuID, count := range userCart {
		items = append(items,
			models.CartItem{
				SkuID: skuID,
				Count: count,
			})
	}

	return items, nil
}
