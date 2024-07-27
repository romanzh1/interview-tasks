package shard_manager

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
)

var (
	ErrShardIndexOutOfRange = errors.New("shard index is out of range")
)

type ShardKey string
type ShardIndex int

type ShardFn func(ShardKey) ShardIndex

type Manager struct {
	fn     ShardFn
	shards []*pgxpool.Pool
}

func GetMurmur3ShardFn(shardsCnt int) ShardFn {
	// ! без seed нам нужен воспроизводимый результат
	hasher := murmur3.New32()
	return func(key ShardKey) ShardIndex {
		defer hasher.Reset()

		_, _ = hasher.Write([]byte(key))

		return ShardIndex(hasher.Sum32() % uint32(shardsCnt))
	}
}

func New(fn ShardFn, shards []*pgxpool.Pool) *Manager {
	return &Manager{
		fn:     fn,
		shards: shards,
	}
}

func (m *Manager) GetShardIndex(key ShardKey) ShardIndex {
	return m.fn(key)
}

func (m *Manager) GetShardIndexFromID(id int64) ShardIndex {
	// 123002
	// 123- seq с шарда
	// 2- номер шарда
	return ShardIndex(id % 1000)
}

func (m *Manager) Pick(index ShardIndex) (*pgxpool.Pool, error) {
	if int(index) < len(m.shards) {
		return m.shards[index], nil
	}

	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, index, len(m.shards))
}
