package main

import (
	"fmt"
)

// 4. Написать что выведется в консоль
// call 2  21
// call 1  30 - копирование переменной в функцию переданную в defer в момент инициализации
// call 0  31 - вызов анонимной функции во время выполнения
func main() {
	a := 10
	defer func() { fmt.Println("call 0 ", a+10) }()
	a = 20
	defer fmt.Println("call 1 ", a+10)

	a++
	fmt.Println("call 2 ", a)
}
