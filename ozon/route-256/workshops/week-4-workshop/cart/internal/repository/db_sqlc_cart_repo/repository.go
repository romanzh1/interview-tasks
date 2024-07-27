package db_sqlc_cart_repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"week-4-workshop/cart/internal/domain"
)

type Storage struct {
	conn *pgx.Conn
	cmd  *Queries
}

func NewStorage(cmd *Queries) *Storage {
	return &Storage{
		cmd: cmd,
	}
}

func (s *Storage) AddItem(ctx context.Context, userID int64, item domain.Item) error {
	err := s.cmd.AddItem(ctx, AddItemParams{
		UserID: userID,
		Sku:    int32(item.SKU),
		Count:  int32(item.Count),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error) {
	items, err := s.cmd.GetItemsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return repackItems(items), nil
}

func repackItems([]Item) []domain.Item {
	return nil
}
