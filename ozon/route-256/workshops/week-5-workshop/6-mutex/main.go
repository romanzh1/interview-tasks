package main

import (
	"fmt"
	"sync"
	"time"
)

type lmap[k comparable, v any] struct {
	mx sync.RWMutex
	m  map[k]v
}

func newLmap[k comparable, v any]() *lmap[k, v] {
	return &lmap[k, v]{
		m: make(map[k]v),
	}
}

func (l *lmap[k, v]) Set(key k, value v) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.m[key] = value
}

func (l *lmap[k, v]) Get(key k) (value v, ok bool) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	value, ok = l.m[key]
	return
}

func main() {
	var (
		m  = newLmap[int, int]()
		wg = sync.WaitGroup{}
	)

	go func() {
		for {
			for i := range 100 {
				v, _ := m.Get(i)
				fmt.Printf("get %v %d\n", v, i)
				if v == 10_000 {
					fmt.Println("never")
				}
			}
		}
	}()

	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Set(i, i)
			fmt.Printf("set %d %d\n", i, i)
		}()
	}

	wg.Wait()
	time.Sleep(time.Second)
}
