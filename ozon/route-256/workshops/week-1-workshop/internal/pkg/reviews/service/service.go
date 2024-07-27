package service

import (
	"context"
	"errors"
	"fmt"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/model"
	"strings"
)

type Repository interface {
	Create(ctx context.Context, r model.Review) (*model.Review, error)
	GetReviews(ctx context.Context, sku int) ([]model.Review, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateReview(ctx context.Context, r model.Review) (*model.Review, error) {
	// validation
	if r.Sku < 0 || len(strings.Trim(r.Comment, " ")) == 0 || r.UserID < 0 {
		return nil, errors.New("review is invalid")
	}

	// save to storage
	updatedReview, err := s.repository.Create(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("repository.Create: %w", err)
	}

	return updatedReview, nil
}

func (s *Service) GetReviewsBySku(ctx context.Context, sku int) ([]model.Review, error) {
	if sku < 1 {
		return nil, errors.New("sku is invalid")
	}

	reviews, err := s.repository.GetReviews(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("repository.GetReviews: %w", err)
	}

	return reviews, nil
}
