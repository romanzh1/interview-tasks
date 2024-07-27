package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/shard_manager"
)

type Repo struct {
	sm shardManager
}

type shardManager interface {
	GetShardIndex(key shard_manager.ShardKey) shard_manager.ShardIndex
	GetShardIndexFromID(id int64) shard_manager.ShardIndex
	Pick(key shard_manager.ShardIndex) (*pgxpool.Pool, error)
}

func NewRepo(sm shardManager) *Repo {
	return &Repo{
		sm: sm,
	}
}
