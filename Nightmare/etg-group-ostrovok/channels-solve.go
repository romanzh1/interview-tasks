package main

import "sync"

// solve
func merge(chs ...chan int) chan int {
	newCh := make(chan int)
	wg := sync.WaitGroup{}

	for _, ch := range chs {
		wg.Add(1)
		go func(ch chan int) {
			for val := range ch {
				newCh <- val
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(newCh)
	}()

	return newCh
}

func main() {
	ch1 := startProducerA()
	ch2 := startProducerB()

	for el := range merge(ch1, ch2) {
		println(el)
	}
}
