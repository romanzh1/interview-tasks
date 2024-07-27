package delete

import (
	"context"
)

type (
	Handler struct{}
)

func New() *Handler {
	return &Handler{}
}

func (h *Handler) DeleteItem(ctx context.Context, userID int64, sku uint32) error {

	return nil
}
