package main

import (
	"fmt"
	"sort"
	"sync"
)

// 1. Починить код, чтобы выводилась сумма всех чисел
func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(v int) {
			defer wg.Done()
			ch <- v
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var sum int
	for v := range ch {
		sum += v
	}

	fmt.Printf("result: %d\n", sum)
}

// 2. Нужно поменять функцию ap так, чтобы 10 попало в выводимый слайс
func main2() {
	v := []int{3, 4, 1, 2, 5}

	ap(&v)
	sr(v)

	fmt.Println(v)
}

func ap(arr *[]int) {
	*arr = append(*arr, 10)
}

func sr(arr []int) {
	sort.Ints(arr)
}

// 3. Что выведет следующий код и почему?
type impl struct{}

type I interface {
	C()
}

func (*impl) C() {}

func A() I {
	return nil
}

func B() I {
	var ret *impl
	return ret
}

// False
// true
// false
// Потому что Go интерфейсное значение состоит из двух частей:
//
//	Тип: указатель на тип значения, хранящегося в интерфейсе.
//	Значение: само значение, которое интерфейс хранит.
//
// Когда функция возвращает интерфейс, возвращаемое значение может быть nil, но тип, связанный с этим интерфейсом, может быть определён.
func main3() {
	a := A()
	b := B()
	fmt.Println(a == b)
	fmt.Println(a == nil)
	fmt.Println(b == nil)
}

// 4. Написать что выведется в консоль
// 11
// 20 - инициализация переменной в defer
// 21 - вызов анонимной функции во время выполнения
func main4() {
	a := 10
	defer func() { fmt.Println("call 0 ", a+10) }()

	defer fmt.Println("call 1 ", a+10)

	a++
	fmt.Println("call 2 ", a)
}
