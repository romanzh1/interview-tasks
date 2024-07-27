package db_cart_repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"week-4-workshop/cart/internal/domain"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(conn *pgx.Conn) *Storage {
	return &Storage{
		conn: conn,
	}
}

func (s *Storage) AddItem(ctx context.Context, userID int64, item domain.Item) error {
	const query = `
	INSERT INTO items(user_id, sku, count)
	VALUES ($1, $2, $3);`

	_, err := s.conn.Exec(ctx, query, userID, item.SKU, item.Count)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) addItem(ctx context.Context, tx pgx.Tx, userID int64, item domain.Item) error {
	const query = `
	INSERT INTO items(user_id, sku, count)
	VALUES ($1, $2, $3);`

	_, err := tx.Exec(ctx, query, userID, item.SKU, item.Count)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddItemInTx(ctx context.Context, userID int64, item domain.Item) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.addItem(ctx, tx, userID, item)
	if err != nil {
		return err
	}

	//s.addLog(ctx, tx)

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
