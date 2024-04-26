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

	return s[999_999:]
}

func printMemStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(m.Alloc / 1024 / 1024)
}
