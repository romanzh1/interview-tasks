// Основная проблема заключается в том, что каждый вызов getSubSlice создает слайс размером 1_000_000 элементов,
// но возвращает только последний элемент. Однако, go всё равно удерживает
// ссылку на весь массив из 1_000_000 элементов, что приводит к значительному потреблению памяти.
package main

import (
	"fmt"
	"runtime"
)

func main() {
	slices()
}

func slices() {
	s := getSubSlice()
	printMemStat()

	all := make([][]int, 0)
	all = append(all, s)

	for i := 1; i < 10; i++ {
		s2 := getSubSlice()
		runtime.GC()
		printMemStat()

		all = append(all, s2)
	}

	runtime.GC()
	printMemStat()

	fmt.Println(all)
}

func getSubSlice() []int {
	s := make([]int, 1_000_000)
	subSlice := make([]int, 1)
	copy(subSlice, s[999_999:])
	return subSlice
}

func printMemStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(m.Alloc / 1024 / 1024)
}
