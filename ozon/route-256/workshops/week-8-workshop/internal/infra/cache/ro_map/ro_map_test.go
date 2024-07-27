package map_mx_repo

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func runGetter(ctx context.Context, m map[int]int, iterCnt int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		// add jitter
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)*100))
		for i := 0; i < iterCnt; i++ {
			if ctx.Err() != nil {
				return
			}

			n := rand.Intn(len(m))
			v, ok := m[n]
			if !ok {
				panic(fmt.Sprintf("not ok for key=%d", n))
			}
			if v != n {
				panic(fmt.Sprintf("v != n : %d != %d for key=%d", v, n, n))
			}
		}
	}()
}

func runWriter(ctx context.Context, m map[int]int, iterCnt int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		// add jitter
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)*100))
		for i := 0; i < iterCnt; i++ {
			if ctx.Err() != nil {
				return
			}

			n := rand.Intn(len(m))
			m[n] = n
		}
	}()
}

// Test_ReadOnlyMap
// go test -race ./internal/infra/cache/ro_map/ro_map_test.go
func Test_ReadOnlyMap(t *testing.T) {
	var (
		cnt        = 10000
		interCnt   = 1000
		workersCnt = 100
		testMap    = make(map[int]int, cnt)
		wg         = &sync.WaitGroup{}
	)

	for i := 0; i < cnt; i++ {
		testMap[i] = i
	}

	for i := 0; i < workersCnt; i++ {
		runGetter(context.Background(), testMap, interCnt, wg)
	}

	// emulate data race
	//runWriter(context.Background(), testMap, interCnt, wg)

	wg.Wait()
}
