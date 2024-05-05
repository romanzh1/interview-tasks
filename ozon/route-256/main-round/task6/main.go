package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func minTransports(n, k int, boxes []int) int {
	weights := make([]int, len(boxes))
	for i, power := range boxes {
		weights[i] = 1 << power //😎
	}

	sort.Slice(weights, func(i, j int) bool {
		return weights[i] > weights[j]
	})

	transports := 0
	for len(weights) > 0 {
		transports++
		currentCapacity := make([]int, n)
		for i := 0; i < len(weights); {
			placed := false

			for j := 0; j < n; j++ {
				if currentCapacity[j]+weights[i] <= k {
					currentCapacity[j] += weights[i]
					weights = append(weights[:i], weights[i+1:]...)
					placed = true
					break
				}
			}

			if !placed {
				i++
			}
		}
	}

	return transports
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	in.Scan()
	NumTests, _ := strconv.Atoi(in.Text())
	results := make([]int, NumTests)

	for i := 0; i < NumTests; i++ {
		in.Scan()
		parts := strings.Fields(in.Text())
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])

		in.Scan()
		m, _ := strconv.Atoi(in.Text())
		boxes := make([]int, m)

		in.Scan()
		boxParts := strings.Fields(in.Text())
		for j := range boxes {
			boxes[j], _ = strconv.Atoi(boxParts[j])
		}

		results[i] = minTransports(n, k, boxes)
	}

	fmt.Fprintln(out, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(results)), "\n"), "[]"))
	out.Flush()
}
