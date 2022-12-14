package main

import "fmt"

// solve  1
func MoveZeros(zeros []int) []int {
	return MoveZero(zeros, 0)
}

func MoveZero(zeros []int, n int) []int {
	if n == len(zeros)-1 {
		return zeros
	} else if zeros[n] == 0 {
		if zeros[n+1] == 0 {
			zeros = MoveZero(zeros, n+1)
		}
		zeros[n], zeros[n+1] = zeros[n+1], zeros[n]
	}

	return MoveZero(zeros, n+1)
}

// solve 2
func Summ(pir [][]int, n int) int {
	sum := 0
	for i := 0; i <= n; i++ {
		sum += pir[n][i]
	}

	return sum
}

func main() {
	fmt.Println(Summ([][]int{{1}, {3, 5}, {7, 9, 11}}, 2))
	fmt.Println(MoveZeros([]int{1, 0, 1, 2, 0, 1, 3}))
}
