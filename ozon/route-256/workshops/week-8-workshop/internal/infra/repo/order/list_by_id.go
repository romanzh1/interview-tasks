package order

import (
	"context"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/shard_manager"
	"golang.org/x/sync/errgroup"
)

func (s *Repo) ListByID(ctx context.Context, orderIDs []int64) ([]domain.Order, error) {

	// orderIDs : [1000, 2001, 3000, 4001]
	// [
	//		0: [1000, 3000],
	//		1: [2001, 4001]
	// ]
	idsByShard := lo.GroupBy(orderIDs, func(orderID int64) shard_manager.ShardIndex {
		return s.sm.GetShardIndexFromID(orderID)
	})

	const query = `
	SELECT id, user_id, description, created_at
	FROM orders
	WHERE id = any($1)`

	errGr := errgroup.Group{}
	resultByShard := make([][]domain.Order, len(idsByShard))

	for shIndex, ids := range idsByShard {
		// не обязательно в gp 1.22
		shIndex, ids := shIndex, ids
		errGr.Go(func() error {
			db, err := s.sm.Pick(shIndex)

			rows, err := db.Query(ctx, query, pq.Array(ids))
			if err != nil {
				return err
			}
			defer rows.Close()

			var (
				rowOrder domain.Order
				orders   []domain.Order
			)

			for rows.Next() {
				err = rows.Scan(&rowOrder.ID, &rowOrder.UserID, &rowOrder.Description, &rowOrder.CreatedAt)
				if err != nil {
					return err
				}

				orders = append(orders, rowOrder)
			}

			if rows.Err() != nil {
				return rows.Err()
			}

			resultByShard[shIndex] = orders

			return nil
		})
	}

	err := errGr.Wait()
	if err != nil {
		return nil, err
	}

	result := make([]domain.Order, 0, len(orderIDs))
	for _, orders := range resultByShard {
		result = append(result, orders...)
	}

	// Может потребоваться дополнительная сортировка, данные
	// будут упорядочены только в рамках шарда

	return result, nil
}
