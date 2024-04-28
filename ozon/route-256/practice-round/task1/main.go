package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	in.Scan()
	n, _ := strconv.Atoi(in.Text())

	sum := make([]int, 0, n)

	for i := 0; i < n; i++ {
		in.Scan()
		nums := strings.Split(in.Text(), " ")

		a, _ := strconv.Atoi(nums[0])
		b, _ := strconv.Atoi(nums[1])
		sum = append(sum, a+b)
	}

	fmt.Fprintln(out, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sum)), "\n"), "[]"))

	out.Flush()
}
