package main

import (
	"fmt"
	"time"
)

func main() {
	numbers := make(chan int)

	go func() {
		v := <-numbers
		time.Sleep(time.Millisecond)
		fmt.Println(v)
	}()

	numbers <- 1
	// numbers <- 2

	time.Sleep(time.Second)
}
