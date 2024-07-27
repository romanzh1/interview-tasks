package order

import (
	"context"
	"net/http"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type orderRepo interface {
	Create(ctx context.Context, order domain.Order) (int64, error)
	GetByOrderID(ctx context.Context, orderID int64) (*domain.Order, error)
	ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error)
	ListByID(ctx context.Context, orderIDs []int64) ([]domain.Order, error)
}

type Handler struct {
	orderRepo orderRepo
}

func New(orderRepo orderRepo) *Handler {
	return &Handler{
		orderRepo: orderRepo,
	}
}

func writeResponse(w http.ResponseWriter, data []byte, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}
