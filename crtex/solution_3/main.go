package main

import (
	"fmt"
)

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
func main() {
	a := A()
	b := B()
	fmt.Println(a == b)
	fmt.Println(a == nil)
	fmt.Println(b == nil)
}
