package order

import (
	"context"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/shard_manager"
	"strconv"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

func (s *Repo) ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	shIndex := s.sm.GetShardIndex(shard_manager.ShardKey(strconv.FormatInt(userID, 10)))
	db, err := s.sm.Pick(shIndex)

	const query = `
	SELECT id, user_id, description, created_at
	FROM orders
	WHERE user_id = $1;`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		rowOrder domain.Order
		orders   []domain.Order
	)

	for rows.Next() {
		err = rows.Scan(&rowOrder.ID, &rowOrder.UserID, &rowOrder.Description, &rowOrder.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, rowOrder)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orders, nil
}
