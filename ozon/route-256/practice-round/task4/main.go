package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Athlete struct {
	Time  int
	Index int
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(calculateString(reader))
}

func calculateString(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	t, _ := strconv.Atoi(strings.TrimSpace(line))

	results := make([]string, t)
	builder := strings.Builder{}

	for i := 0; i < t; i++ {
		line, _ = reader.ReadString('\n')
		n, _ := strconv.Atoi(strings.TrimSpace(line))

		times := make([]Athlete, n)
		line, _ = reader.ReadString('\n')
		splits := strings.Fields(line)

		for j := 0; j < n; j++ {
			time, _ := strconv.Atoi(splits[j])
			times[j] = Athlete{Time: time, Index: j}
		}

		// Сортировка по времени
		sort.Slice(times, func(a, b int) bool {
			return times[a].Time < times[b].Time
		})

		ranks := make([]int, n)
		currentRank := 1

		for j := 0; j < n; j++ {
			if j > 0 && times[j].Time-times[j-1].Time > 1 {
				currentRank = j + 1
			}

			ranks[times[j].Index] = currentRank
		}

		for _, rank := range ranks {
			builder.WriteString(strconv.Itoa(rank) + " ")
		}
		builder.WriteString("\n")
		results[i] = strings.TrimSpace(builder.String())
		builder.Reset()
	}

	return strings.Join(results, "\n") + "\n"
}
