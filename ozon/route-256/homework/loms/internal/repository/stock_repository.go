package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"route256/libs/metrics"

	"route256/loms/internal/models"
)

type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) GetBySKU(ctx context.Context, sku uint32) (models.Stock, error) {
	start := time.Now()
	var stock models.Stock
	err := r.db.QueryRow(ctx, `
		SELECT sku, total_count, reserved
		FROM stocks
		WHERE sku = $1`, sku).Scan(&stock.SKU, &stock.TotalCount, &stock.Reserved)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	metrics.DBRequests.WithLabelValues("select", status).Inc()
	metrics.DBRequestDuration.WithLabelValues("select", status).Observe(duration)

	return stock, err
}

func (r *StockRepository) ReserveStocks(ctx context.Context, tx pgx.Tx, items []models.OrderItem) error {
	for _, item := range items {
		start := time.Now()
		_, err := tx.Exec(ctx, `
			UPDATE stocks
			SET reserved = reserved + $2
			WHERE sku = $1 AND total_count >= reserved + $2`, item.SkuID, item.Quantity)
		duration := time.Since(start).Seconds()

		status := "success"
		if err != nil {
			status = "error"
		}

		metrics.DBRequests.WithLabelValues("update", status).Inc()
		metrics.DBRequestDuration.WithLabelValues("update", status).Observe(duration)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *StockRepository) RemoveReservedStocks(ctx context.Context, tx pgx.Tx, items []models.OrderItem) error {
	for _, item := range items {
		start := time.Now()
		_, err := tx.Exec(ctx, `
			UPDATE stocks
			SET reserved = reserved - $2
			WHERE sku = $1 AND reserved >= $2`, item.SkuID, item.Quantity)
		duration := time.Since(start).Seconds()

		status := "success"
		if err != nil {
			status = "error"
		}

		metrics.DBRequests.WithLabelValues("update", status).Inc()
		metrics.DBRequestDuration.WithLabelValues("update", status).Observe(duration)

		if err != nil {
			return err
		}
	}
	return nil
}
