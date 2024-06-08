package repository

import (
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

func (r *Repository) AddItemToUserCart(userID, skuID int64, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.carts[userID] == nil {
		r.carts[userID] = make(map[int64]uint16)
	}

	r.carts[userID][skuID] += count

	return nil
}

func (r *Repository) DeleteItemFromUserCart(userID, skuID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts[userID], skuID)

	return nil
}

func (r *Repository) ClearUserCart(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts, userID)

	return nil
}

func (r *Repository) ListUserCart(userID int64) ([]models.CartItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userCart, exists := r.carts[userID]
	if !exists || len(userCart) == 0 {
		return nil, models.CartIsEmpty
	}

	var items []models.CartItem
	for skuID, count := range userCart {
		items = append(items,
			models.CartItem{
				SkuID: skuID,
				Count: count,
			})
	}

	return items, nil
}
