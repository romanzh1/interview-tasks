package list

import (
	"context"
	"week-4-workshop/cart/internal/domain"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) List(ctx context.Context, userID int64) ([]domain.ListItem, error) {

	return []domain.ListItem{}, nil
}
