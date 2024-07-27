package order

import (
	"context"
	"strconv"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/shard_manager"
)

func (s *Repo) Create(ctx context.Context, order domain.Order) (int64, error) {
	shIndex := s.sm.GetShardIndex(shard_manager.ShardKey(strconv.FormatInt(order.UserID, 10)))
	db, err := s.sm.Pick(shIndex)
	if err != nil {
		return 0, err
	}

	// Варианты как можно раскладывать данные по шардам
	// 1. Определяем случайно любой шард
	// shIndex := rand.Int()
	// 2. Определяем например по userID
	// данный пример, заказы одного пользователя всегда на одном шарде
	// эфективное чение списка закащов пользователя
	// неравномерное распределение данных по шардам
	// 3. Можем подмешать номер шарда в идентификатор (что будет если кол-во шардов изменится)

	var orderID int64

	const query = `
	INSERT INTO orders(id, user_id, description, created_at)
	VALUES (nextval('order_id_manual_seq') + $1, $2, $3, $4)
	RETURNING id;`

	// order.UserID
	// murmur3(order.UserID) % shard_cnt = 1 // shardIndex
	// shards[shardIndex]
	// order_id_manual_seq = 42000
	// nextval('order_id_manual_seq') -> 43000
	// 43 + shardIndex + 1000 = 43001

	err = db.QueryRow(ctx, query, shIndex, order.UserID, order.Description, order.CreatedAt).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
