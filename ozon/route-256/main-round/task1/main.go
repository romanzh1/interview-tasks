package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(calculateString(reader))
}

func calculateString(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	num := strings.Fields(line)
	first, _ := strconv.Atoi(num[0])
	second, _ := strconv.Atoi(num[1])

	return strconv.Itoa(first - second)
}
