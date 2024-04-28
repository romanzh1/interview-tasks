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

	scanner.Scan() //number of sets
	firstString := scanner.Text()
	numberOfSets, _ := strconv.Atoi(firstString)
	result := ``

	for in := 0; in < numberOfSets; in++ {
		scanner.Scan() //length of strings: number of strings, blue words, red words, black word
		secondString := strings.Fields(scanner.Text())

		numberOfWords, _ := strconv.Atoi(secondString[0])
		blue, _ := strconv.Atoi(secondString[1])
		red, _ := strconv.Atoi(secondString[2])
		black, _ := strconv.Atoi(secondString[3])

		words := make([]string, 0, numberOfWords)
		shortestBlueWordI := 0
		for i := 0; i < numberOfWords; i++ {
			scanner.Scan()
			words = append(words, scanner.Text())

			if i < blue && len(words[i]) < len(words[shortestBlueWordI]) {
				shortestBlueWordI = i
			}
		}

		shortestBlueWord := []rune(words[shortestBlueWordI])
		delta := 0
		finalWord := ""

		for left, right := 0, len(shortestBlueWord)-1; left != len(shortestBlueWord)-1; right-- {
			if subStr, ok := findSubstr(string(shortestBlueWord), words[black-1], left, right); ok {
				numberOfBlues := checkBlue(words[:blue], subStr)
				numberOfReds := checkRed(words[blue:blue+red], subStr)

				if numberOfBlues-numberOfReds > delta {
					delta = numberOfBlues - numberOfReds
					finalWord = subStr
				}

				fmt.Printf("left = %d, right = %d, shortest string: %s, subStr = %s, ok = %t\n",
					left, right, string(shortestBlueWord), subStr, ok)
			}

			if right == left {
				right = len(shortestBlueWord) + 1
				left++
				continue
			}
		}

		if delta != 0 {
			result += finalWord + " " + strconv.Itoa(delta)
		} else {
			result += words[shortestBlueWordI] + " " + strconv.Itoa(0)
		}

		if in != numberOfSets-1 {
			result += "\n"
		}
	}

	fmt.Println(result)
}

func findSubstr(word, blackWord string, l, r int) (string, bool) {
	currentString := []rune(word)

	if strings.Contains(word, string(currentString[l:r])) &&
		checkBlack(blackWord, currentString[l:r]) && l+r != len(currentString) {
		return string(currentString[l:r]), true
	}

	return "", false
}

func checkBlack(blackWord string, currentString []rune) bool {
	return !strings.Contains(blackWord, string(currentString))
}

func checkBlue(words []string, substr string) int {
	countOccurrences := 0

	for i := 0; i < len(words); i++ {
		if strings.Contains(words[i], substr) {
			countOccurrences++

			if strings.EqualFold(words[i], substr) {
				return 0
			}
		}
	}

	return countOccurrences
}

func checkRed(words []string, substr string) int {
	countOccurrences := 0

	for i := 0; i < len(words); i++ {
		if strings.Contains(words[i], substr) {
			countOccurrences++
		}
	}

	return countOccurrences
}
