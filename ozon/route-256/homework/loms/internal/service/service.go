package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"route256/loms/internal/models"
)

type cartRepository interface {
	CreateOrder(ctx context.Context, tx pgx.Tx, userID int64, status models.OrderStatus) (int64, error)
	GetOrder(ctx context.Context, orderID int64) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, tx pgx.Tx, orderID int64, status models.OrderStatus) error
	CancelOrder(ctx context.Context, tx pgx.Tx, orderID int64) error
}

type stockRepository interface {
	ReserveStocks(ctx context.Context, tx pgx.Tx, items []models.OrderItem) error
	GetBySKU(ctx context.Context, sku uint32) (models.Stock, error)
	RemoveReservedStocks(ctx context.Context, tx pgx.Tx, items []models.OrderItem) error
}

type TxManager interface {
	Tx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error, opts *pgx.TxOptions) error
}

type Service struct {
	txManager TxManager
	orderRepo cartRepository
	stockRepo stockRepository
}

func NewService(txManager TxManager, orderRepo cartRepository, stockRepo stockRepository) *Service {
	return &Service{txManager: txManager, orderRepo: orderRepo, stockRepo: stockRepo}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, items []models.OrderItem) (int64, error) {
	var orderID int64

	fmt.Println(items, userID)
	err := s.txManager.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var err error

		orderID, err = s.orderRepo.CreateOrder(ctx, tx, userID, models.OrderStatusNew)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		err = s.stockRepo.ReserveStocks(ctx, tx, items)
		if err != nil {
			_ = s.orderRepo.UpdateOrderStatus(ctx, tx, orderID, models.OrderStatusFailed)
			return fmt.Errorf("failed to reserve stocks: %w", err)
		}

		err = s.orderRepo.UpdateOrderStatus(ctx, tx, orderID, models.OrderStatusAwaitingPayment)
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	}, nil)

	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (s *Service) GetOrder(ctx context.Context, orderID int64) (models.Order, error) {
	return s.orderRepo.GetOrder(ctx, orderID)
}

func (s *Service) CancelOrder(ctx context.Context, orderID int64) error {
	return s.txManager.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err := s.orderRepo.UpdateOrderStatus(ctx, tx, orderID, models.OrderStatusCanceled)
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return s.orderRepo.CancelOrder(ctx, tx, orderID)
	}, nil)
}

func (s *Service) GetStockInfo(ctx context.Context, sku uint32) (uint64, error) {
	var available uint64

	err := s.txManager.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		stock, err := s.stockRepo.GetBySKU(ctx, sku)
		if err != nil {
			return err
		}

		available = stock.TotalCount - stock.Reserved
		return nil
	}, nil)

	return available, err
}

func (s *Service) OrderPay(ctx context.Context, orderID int64) error {
	return s.txManager.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		order, err := s.orderRepo.GetOrder(ctx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}

		if order.Status != models.OrderStatusAwaitingPayment {
			return fmt.Errorf("order is not in awaiting payment status")
		}

		err = s.stockRepo.RemoveReservedStocks(ctx, tx, order.Items)
		if err != nil {
			return fmt.Errorf("failed to remove reserved stocks: %w", err)
		}

		err = s.orderRepo.UpdateOrderStatus(ctx, tx, orderID, models.OrderStatusPayed)
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	}, nil)
}
