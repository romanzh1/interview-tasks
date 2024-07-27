package map_mx_repo

import (
	"sync"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type Repo struct {
	data       map[int64]*domain.Order
	maxOrderID int64
	mx         sync.Mutex
}

func New(orders []domain.Order) *Repo {
	var (
		data       = make(map[int64]*domain.Order, len(orders))
		maxOrderID = int64(0)
	)

	for i := 0; i < len(orders); i++ {
		data[orders[i].ID] = &orders[i]
	}

	return &Repo{
		data:       data,
		maxOrderID: maxOrderID,
	}
}

func (r *Repo) Get(orderID int64) *domain.Order {
	r.mx.Lock()
	defer r.mx.Unlock()

	if order, ok := r.data[orderID]; ok {
		return order
	}

	return nil
}

func (r *Repo) Add(order *domain.Order) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.data[order.ID] = order
}
