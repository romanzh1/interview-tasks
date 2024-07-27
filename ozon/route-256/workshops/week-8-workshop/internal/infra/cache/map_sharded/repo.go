package map_mx_repo

import (
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/cache/map_mx_repo"
)

type Repo struct {
	data      []*map_mx_repo.Repo
	shardsCnt int64
}

func New(orders []domain.Order, shardsCnt int64) *Repo {
	var (
		data = make([]*map_mx_repo.Repo, shardsCnt)
	)

	// обязательно инициализируем map
	for i := 0; i < int(shardsCnt); i++ {
		data[i] = map_mx_repo.New([]domain.Order{})
	}

	var shardID int64
	for i := 0; i < len(orders); i++ {
		// функция шардирования
		shardID = orders[i].ID % shardsCnt

		data[shardID].Add(&orders[i])
	}

	return &Repo{
		data:      data,
		shardsCnt: shardsCnt,
	}
}

func (r *Repo) Get(orderID int64) *domain.Order {
	// не нужны, блокировки на уровне шардов
	//r.mx.Lock()
	//defer r.mx.Unlock()

	if order := r.data[orderID%r.shardsCnt].Get(orderID); order != nil {
		return order
	}

	return nil
}

func (r *Repo) Add(order *domain.Order) {
	// не нужны, блокировки на уровне шардов
	//r.mx.Lock()
	//defer r.mx.Unlock()

	r.data[order.ID%r.shardsCnt].Add(order)
}
