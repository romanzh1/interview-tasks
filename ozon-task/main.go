package main

import "time"

// func a() {
// 	x := []int{}
// 	x = append(x, 0)
// 	x = append(x, 1)
// 	x = append(x, 2)
// 	fmt.Println(cap(x), len(x))
// 	y := append(x, 3)
// 	fmt.Println(cap(x), len(x))

// 	z := append(x, 4)
// 	fmt.Println(cap(x), len(x))

// 	fmt.Println(x, y, z)
// }

// func main() {
// 	a()
// }

func main() {
	timeStart := time.Now()
	_, _ = <-worker(), <-worker()
	println(int(time.Since(timeStart).Seconds()))
}

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 1
	}()
	return ch
}
