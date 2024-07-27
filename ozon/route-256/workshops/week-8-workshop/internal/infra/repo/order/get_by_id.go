package order

import (
	"context"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

func (s *Repo) GetByOrderID(ctx context.Context, orderID int64) (*domain.Order, error) {
	shIndex := s.sm.GetShardIndexFromID(orderID)
	db, err := s.sm.Pick(shIndex)

	const query = `
	SELECT user_id, description, created_at
	FROM orders
	WHERE id = $1;`

	rows, err := db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := &domain.Order{
		ID: orderID,
	}

	for rows.Next() {
		err = rows.Scan(&order.UserID, &order.Description, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		break
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return order, nil
}
