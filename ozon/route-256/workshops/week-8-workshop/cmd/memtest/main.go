package main

import (
	"fmt"
	"runtime"
	"time"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	PrintMemUsage()
	fmt.Println("begin")
	ar := func() []int {
		ar := make([]int, 1000000)
		//return ar[5000:5001]
		_ = ar
		return []int{}
	}()
	fmt.Println("end")
	PrintMemUsage()
	time.Sleep(5)
	fmt.Println("after sleep")
	// Force GC to clear up, should see a memory drop
	runtime.GC()
	PrintMemUsage()
	fmt.Println(len(ar))
}
