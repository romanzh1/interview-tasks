package map_mx_repo

import (
	"math/rand"
	"testing"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
	"golang.org/x/sync/errgroup"
)

func initRepo(factory *domain.OrderFactory, ordersCount int) *Repo {
	orders := make([]domain.Order, ordersCount)
	for i := 0; i < ordersCount; i++ {
		orders[i] = factory.Create()
	}

	return New(orders)
}

// Benchmark_ReadOnly-8   	66951518	        17.86 ns/op
// Benchmark_ReadOnly-8   	67765548	        17.56 ns/op // rw mx
// Benchmark_ReadOnly-8   	71069871	        16.79 ns/op // sharded
// Benchmark_ReadOnly-8   	79563504	        15.30 ns/op // sync map
func Benchmark_ReadOnly(b *testing.B) {
	var (
		ordersCount = 100000
		factory     = domain.NewDefaultFactory()
		repo        = initRepo(factory, ordersCount)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.Get(rand.Int63n(int64(ordersCount)))
	}
}

// Benchmark_ConcurrentReadOnly-8   	 4328586	       272.8 ns/op // mx
// Benchmark_ConcurrentReadOnly-8   	 4279719	       279.2 ns/op // rw mx
// Benchmark_ConcurrentReadOnly-8   	 4452704	       269.1 ns/op // sharded
// Benchmark_ConcurrentReadOnly-8   	 4387389	       272.4 ns/op // sync map
func Benchmark_ConcurrentReadOnly(b *testing.B) {
	var (
		ordersCount = 100000
		factory     = domain.NewDefaultFactory()
		repo        = initRepo(factory, ordersCount)
	)

	eg := errgroup.Group{}
	eg.SetLimit(10)

	fn := func() error {
		_ = repo.Get(rand.Int63n(int64(ordersCount)))
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eg.Go(fn)
	}
}

// Benchmark_ConcurrentReadWrite-8   	 3940653	       316.5 ns/op // mx
// Benchmark_ConcurrentReadWrite-8   	 3784472	       312.2 ns/op // rw mx
// Benchmark_ConcurrentReadWrite-8   	 3834913	       304.4 ns/op // sharded
// Benchmark_ConcurrentReadWrite-8   	 4025670	       298.9 ns/op // sync map
func Benchmark_ConcurrentReadWrite(b *testing.B) {
	var (
		ordersCount = 100000
		factory     = domain.NewDefaultFactory()
		repo        = initRepo(factory, ordersCount)
	)

	eg := errgroup.Group{}
	eg.SetLimit(10)

	fnRead := func() error {
		_ = repo.Get(rand.Int63n(int64(ordersCount)))
		return nil
	}

	fnWrite := func() error {
		order := factory.Create()
		repo.Add(&order)
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 == 1 {
			eg.Go(fnWrite)
		} else {
			eg.Go(fnRead)
		}
	}
}

// Benchmark_ConcurrentReadWrite1000-8   	 2983358	       406.3 ns/op // mx
// Benchmark_ConcurrentReadWrite1000-8   	 3038359	       386.7 ns/op // rw mx
// Benchmark_ConcurrentReadWrite1000-8   	 3876193	       310.5 ns/op // sharded
// Benchmark_ConcurrentReadWrite1000-8   	 3712120	       303.5 ns/op // sync map
func Benchmark_ConcurrentReadWrite1000(b *testing.B) {
	var (
		ordersCount = 100000
		factory     = domain.NewDefaultFactory()
		repo        = initRepo(factory, ordersCount)
	)

	eg := errgroup.Group{}
	eg.SetLimit(1000)

	fnRead := func() error {
		_ = repo.Get(rand.Int63n(int64(ordersCount)))
		return nil
	}

	fnWrite := func() error {
		order := factory.Create()
		repo.Add(&order)
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 == 1 {
			eg.Go(fnWrite)
		} else {
			eg.Go(fnRead)
		}
	}
}
