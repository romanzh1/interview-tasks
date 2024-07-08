package repository

import (
	"context"
	"sync"
	"time"

	"route256/cart/internal/models"
	"route256/libs/metrics"
)

type Repository struct {
	carts map[int64]map[int64]uint16 // userID -> skuID -> count
	mu    sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{carts: make(map[int64]map[int64]uint16)}
}

func (r *Repository) AddItemToUserCart(ctx context.Context, cart models.CartRequest) error {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.carts[cart.UserID] == nil {
		r.carts[cart.UserID] = make(map[int64]uint16)
	}

	r.carts[cart.UserID][cart.SkuID] += cart.Count

	duration := time.Since(start).Seconds()
	statusLabel := "success"

	metrics.DBRequests.WithLabelValues("create", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("create", statusLabel).Observe(duration)

	return nil
}

func (r *Repository) DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts[userID], skuID)

	duration := time.Since(start).Seconds()
	statusLabel := "success"

	metrics.DBRequests.WithLabelValues("delete", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("delete", statusLabel).Observe(duration)

	return nil
}

func (r *Repository) ClearUserCart(ctx context.Context, userID int64) error {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts, userID)

	duration := time.Since(start).Seconds()
	statusLabel := "success"

	metrics.DBRequests.WithLabelValues("delete", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("delete", statusLabel).Observe(duration)

	return nil
}

func (r *Repository) ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, error) {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	userCart, exists := r.carts[userID]
	if !exists || len(userCart) == 0 {
		return nil, nil
	}

	items := make([]models.CartItem, 0, len(userCart))
	for skuID, count := range userCart {
		items = append(items,
			models.CartItem{
				SkuID: skuID,
				Count: count,
			})
	}
	statusLabel := "success"

	duration := time.Since(start).Seconds()

	metrics.DBRequests.WithLabelValues("select", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("select", statusLabel).Observe(duration)

	return items, nil
}
