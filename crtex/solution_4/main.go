package main

import (
	"fmt"
)

// 4. Написать что выведется в консоль
// 11
// 20 - инициализация переменной в defer
// 21 - вызов анонимной функции во время выполнения
func main() {
	a := 10
	defer func() { fmt.Println("call 0 ", a+10) }()

	defer fmt.Println("call 1 ", a+10)

	a++
	fmt.Println("call 2 ", a)
}
