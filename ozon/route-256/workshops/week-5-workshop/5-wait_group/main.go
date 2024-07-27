package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const counter int64 = 10000

func main() {
	res := iterateAtomic()
	fmt.Println(res)
}

func iterate() int64 {
	var a int64
	for i := range counter {
		a += i
	}
	return a
}

func iterateInt() int64 {
	wg := sync.WaitGroup{}
	wg.Add(int(counter))

	var a int64
	for i := range counter {
		go func() {
			defer wg.Done()
			a += i
		}()
	}

	wg.Wait()
	return a
}

func iterateAtomic() int64 {
	wg := sync.WaitGroup{}
	wg.Add(int(counter))
	a := atomic.Int64{}
	for i := range counter {
		go func() {
			defer wg.Done()
			a.Add(i)
		}()
	}
	wg.Wait()
	return a.Load()
}

func iterateMutex() int64 {
	wg := sync.WaitGroup{}
	wg.Add(int(counter))
	mx := sync.Mutex{}

	var a int64
	for i := range counter {
		go func() {
			defer func() {
				mx.Unlock()
				wg.Done()
			}()
			mx.Lock()
			a += i
		}()
	}
	wg.Wait()
	return a
}
