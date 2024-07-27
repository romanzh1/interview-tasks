package main

import (
	"testing"
)

// go test --bench=.

// 3331
func Benchmark_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterate()
	}
}

// 1911095
func Benchmark_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateInt()
	}
}

// 1853492
func Benchmark_Atomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateAtomic()
	}
}

// 1957542
func Benchmark_Mutext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateMutex()
	}
}
