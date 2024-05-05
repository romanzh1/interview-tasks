package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func compareMaps(a, b map[rune]int) bool {
	for key, val := range a {
		if b[key] != val {
			return false
		}
	}

	return true
}

func createCharMap(s string) map[rune]int {
	charMap := make(map[rune]int)
	for _, char := range s {
		charMap[char]++
	}

	return charMap
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	in.Scan()
	input := strings.Fields(in.Text())
	t, _ := strconv.Atoi(input[1])

	in.Scan()
	letters := in.Text()

	// Карта частот для строки с буквами
	letterMap := createCharMap(letters)

	result := ""
	for i := 0; i < t; i++ {
		in.Scan()
		pin := in.Text()
		pinMap := createCharMap(pin)

		if compareMaps(pinMap, letterMap) {
			result += "YES\n"
		} else {
			result += "NO\n"
		}
	}

	fmt.Fprint(out, result)
	out.Flush()
}
