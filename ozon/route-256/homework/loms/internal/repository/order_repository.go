package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"route256/libs/metrics"

	"route256/loms/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, tx pgx.Tx, userID int64, status models.OrderStatus) (int64, error) {
	start := time.Now()
	var orderID int64
	err := tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, status)
		VALUES ($1, $2)
		RETURNING id`, userID, status).Scan(&orderID)
	duration := time.Since(start).Seconds()

	statusLabel := "success"
	if err != nil {
		statusLabel = "error"
	}

	metrics.DBRequests.WithLabelValues("create", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("create", statusLabel).Observe(duration)

	return orderID, err
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, tx pgx.Tx, orderID int64, status models.OrderStatus) error {
	start := time.Now()
	_, err := tx.Exec(ctx, `
		UPDATE orders
		SET status = $2, updated_at = NOW()
		WHERE id = $1`, orderID, status)
	duration := time.Since(start).Seconds()

	statusLabel := "success"
	if err != nil {
		statusLabel = "error"
	}

	metrics.DBRequests.WithLabelValues("update", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("update", statusLabel).Observe(duration)

	return err
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderID int64) (models.Order, error) {
	start := time.Now()
	var order models.Order
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, status
		FROM orders
		WHERE id = $1`, orderID).Scan(&order.ID, &order.UserID, &order.Status)
	duration := time.Since(start).Seconds()

	statusLabel := "success"
	if err != nil {
		statusLabel = "error"
	}

	metrics.DBRequests.WithLabelValues("select", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("select", statusLabel).Observe(duration)

	return order, err
}

func (r *OrderRepository) CancelOrder(ctx context.Context, tx pgx.Tx, orderID int64) error {
	start := time.Now()
	_, err := tx.Exec(ctx, `
		DELETE FROM orders
		WHERE id = $1`, orderID)
	duration := time.Since(start).Seconds()

	statusLabel := "success"
	if err != nil {
		statusLabel = "error"
	}

	metrics.DBRequests.WithLabelValues("delete", statusLabel).Inc()
	metrics.DBRequestDuration.WithLabelValues("delete", statusLabel).Observe(duration)

	return err
}
