package main

import (
	"fmt"
)

func join(meet [][]int) [][]int {
	intervals := make([][]int, 0)

	for i := 0; i < len(meet)-1; i++ {
		if meet[i][1] > meet[i+1][0] {
			intervals = append(intervals, []int{meet[i][0], meet[i+1][1]})
		} else {
			intervals = append(intervals, meet[i])
		}
	}

	return intervals
}

func main() {
	input := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("input: %v\n", input)
	fmt.Printf("output: %v\n", join(input))
}
