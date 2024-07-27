package map_mx_repo

import (
	"sync"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type Repo struct {
	data *sync.Map
}

func New(orders []domain.Order) *Repo {
	var (
		d = &sync.Map{}
	)
	for i := 0; i < len(orders); i++ {
		d.Store(orders[i].ID, &orders[i])
	}

	return &Repo{
		data: d,
	}
}

func (r *Repo) Get(orderID int64) *domain.Order {
	order, ok := r.data.Load(orderID)
	if !ok || order == nil {
		return nil
	}

	return order.(*domain.Order)
}

func (r *Repo) Add(order *domain.Order) {
	r.data.Store(order.ID, order)
}
