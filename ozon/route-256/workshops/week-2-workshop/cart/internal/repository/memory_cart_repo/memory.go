package memory_cart_repo

import (
	"context"
	"sync"

	"week-2-workshop/cart/internal/domain"
)

// itemsMap is index sku -> item.
type (
	itemsMap map[uint32]domain.Item

	MemoryStorage struct {
		items map[int64]itemsMap
		mtx   sync.RWMutex
	}
)

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[int64]itemsMap),
		mtx:   sync.RWMutex{},
	}
}

func (m *MemoryStorage) AddItem(_ context.Context, userID int64, item domain.Item) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.items[userID] == nil {
		m.items[userID] = itemsMap{}
	}

	m.items[userID][item.SKU] = domain.Item{
		SKU:   item.SKU,
		Count: m.items[userID][item.SKU].Count + item.Count,
	}
	//for key, val := range m.items[userID] {
	//	log.Printf("%v:%v\n", key, val)
	//}
	return nil
}
