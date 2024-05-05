package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Truck struct {
	start, end, capacity, index int
}

type Order struct {
	time, index int
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	in.Scan()
	t, _ := strconv.Atoi(in.Text())

	for ; t > 0; t-- {
		in.Scan()
		n, _ := strconv.Atoi(in.Text())

		in.Scan()
		arrivalStr := strings.Split(in.Text(), " ")
		orders := make([]Order, n)
		for i, v := range arrivalStr {
			time, _ := strconv.Atoi(v)
			orders[i] = Order{time, i}
		}

		in.Scan()
		m, _ := strconv.Atoi(in.Text())

		trucks := make([]Truck, m)
		for j := 0; j < m; j++ {
			in.Scan()
			truckData := strings.Split(in.Text(), " ")
			start, _ := strconv.Atoi(truckData[0])
			end, _ := strconv.Atoi(truckData[1])
			capacity, _ := strconv.Atoi(truckData[2])
			trucks[j] = Truck{start, end, capacity, j}
		}

		sort.Slice(orders, func(i, j int) bool {
			return orders[i].time < orders[j].time
		})

		results := make([]int, n)
		for _, order := range orders {
			found := false

			sort.Slice(trucks, func(i, j int) bool {
				if trucks[i].start == trucks[j].start {
					return trucks[i].index < trucks[j].index
				}
				return trucks[i].start < trucks[j].start
			})

			for j := 0; j < m && !found; j++ {
				if order.time >= trucks[j].start && order.time <= trucks[j].end && trucks[j].capacity > 0 {
					results[order.index] = trucks[j].index + 1
					trucks[j].capacity--
					found = true
				}
			}

			if !found {
				results[order.index] = -1
			}
		}

		for _, result := range results {
			fmt.Fprintf(out, "%d ", result)
		}

		fmt.Fprintln(out)
	}
}
