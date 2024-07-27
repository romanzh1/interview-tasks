package main

import (
	"fmt"
	"time"
)

func main() {
	// simpleGenerator()
	// simpleGeneratorWithClose()
	generator()
}

func simpleGenerator() {
	numbers := make(chan int)

	go func() {
		for i := range 100 {
			numbers <- i
		}
	}()

	for i := range numbers {
		fmt.Println(i)
		time.Sleep(time.Millisecond * 200)
	}
}

func simpleGeneratorWithClose() {
	numbers := make(chan int)

	go func() {
		defer close(numbers)

		for i := range 100 {
			numbers <- i
		}
	}()

	//close(numbers)
	//time.Sleep(time.Millisecond * 200)

	for i := range numbers {
		fmt.Println(i)
	}
}

func generator() {
	generator := func() <-chan int {
		// благодаря буферу, мы раньше завершаем запись и освобождаем ресурсы
		numbers := make(chan int, 20)

		go func() {
			defer func() {
				close(numbers)
				fmt.Println("CLOSE IT!")
			}()
			for i := range 100 {
				numbers <- i
			}
		}()

		return numbers
	}

	numbers := generator()

	//close(numbers)
	//time.Sleep(time.Millisecond * 200)

	for i := range numbers {
		time.Sleep(time.Millisecond * 200)
		fmt.Println(i)
	}
}
