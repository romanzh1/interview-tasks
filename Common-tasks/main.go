package main

import (
	"fmt"
)

func main() {
	// 1
	// Вывести длину строки в символах и количество байт
	str := "Привет, мир!"
	length := len(str)
	fmt.Println("Длина строки в символах:", length)

	length = len([]byte(str))
	fmt.Println("Длина строки в байтах:", length)

	// 2
	// Заменить символ в строке
	str = "Hello, world!"
	oldChar := 'l'
	newChar := 'X'

	newStr := replaceChar(str, oldChar, newChar)
	fmt.Println(newStr)
}

func replaceChar(input string, oldChar rune, newChar rune) string {
	bytes := []byte(input)

	for i, b := range bytes {
		if rune(b) == oldChar {
			bytes[i] = byte(newChar)
		}
	}

	return string(bytes)
}
