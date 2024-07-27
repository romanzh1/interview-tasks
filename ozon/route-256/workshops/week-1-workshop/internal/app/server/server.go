package server

import (
	"context"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/model"
)

type ReviewsService interface {
	CreateReview(ctx context.Context, r model.Review) (*model.Review, error)
	GetReviewsBySku(ctx context.Context, sku int) ([]model.Review, error)
}

type Server struct {
	reviewsService ReviewsService
}

func NewServer(
	reviewsService ReviewsService,
) *Server {
	return &Server{reviewsService: reviewsService}
}
