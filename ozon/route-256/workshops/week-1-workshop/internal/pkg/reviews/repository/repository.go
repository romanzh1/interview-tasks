package repository

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/model"
)

type Storage map[int][]model.Review

type InMemoryRepository struct {
	storage Storage
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(Storage, 10),
	}
}

// TODO context
func (r *InMemoryRepository) Create(_ context.Context, review model.Review) (*model.Review, error) {
	// validation

	r.storage[review.Sku] = append(r.storage[review.Sku], review)

	return &review, nil
}

// TODO context
func (r *InMemoryRepository) GetReviews(_ context.Context, sku int) ([]model.Review, error) {
	if reviews, ok := r.storage[sku]; ok {
		if len(reviews) == 0 {
			return nil, fmt.Errorf("reviews not found by sku %d", sku)
		}
		return reviews, nil
	}

	// TODO errors.Is
	return nil, fmt.Errorf("reviews not found by sku %d", sku)
}
