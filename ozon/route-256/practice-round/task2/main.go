package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	firstString := scanner.Text()

	scanner.Scan()
	secondString := scanner.Text()

	n, _ := strconv.Atoi(secondString)

	originalStickers := []rune(firstString)

	for i := 0; i < n; i++ {
		scanner.Scan()
		parts := strings.Fields(scanner.Text())

		start, _ := strconv.Atoi(parts[0])
		r := parts[2]

		for j, ch := range r {
			originalStickers[start-1+j] = ch
		}
	}

	fmt.Print(string(originalStickers))
}
