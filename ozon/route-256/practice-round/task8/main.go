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
	result := ""

	for in := 0; in < numberOfSets; in++ {
		scanner.Scan() //length of strings: number of strings, blue words, red words, black word
		secondString := strings.Fields(scanner.Text())

		numberOfWords, _ := strconv.Atoi(secondString[0])
		blue, _ := strconv.Atoi(secondString[1])
		red, _ := strconv.Atoi(secondString[2])
		black, _ := strconv.Atoi(secondString[3])

		words := make([]string, 0, numberOfWords)
		wordsMap := make(map[string]bool)
		shortestBlueWordI := 0
		for i := 0; i < numberOfWords; i++ {
			scanner.Scan()
			words = append(words, scanner.Text())
			wordsMap[words[i]] = true

			if i < blue && len(words[i]) < len(words[shortestBlueWordI]) {
				shortestBlueWordI = i
			}
		}

		shortestBlueWord := []rune(words[shortestBlueWordI])
		delta := 0
		finalWord := ""

		// TODO fix bug (где-то насчитал 0, а там 1), для этого надо написать тесты и прогнать все тест кейсы
		// TODO не искать другие слова, если использовал минимальное из синих и оно подошло?
		// TODO в любое случае оптимизировать перебор строк для рассмотрения: брать минимально возможные по очереди

		for jb := 0; jb <= blue; jb++ {
			for left, right := 0, len(shortestBlueWord)-1; left != len(shortestBlueWord)-1; right-- {
				if subStr, ok := findSubstr(string(shortestBlueWord), words[black-1], left, right); ok {
					numberOfBlues := countOccurrences(words[:blue], subStr)
					numberOfReds := countOccurrences(words[blue:blue+red], subStr)

					if numberOfBlues-numberOfReds > delta {
						if !wordsMap[subStr] {
							delta = numberOfBlues - numberOfReds
							finalWord = subStr
						}
					}
				}

				if right == left {
					right = len(shortestBlueWord) + 1
					left++
					continue
				}
			}

			if jb != shortestBlueWordI {
				shortestBlueWord = []rune(words[jb])
			}
		}

		if delta != 0 {
			result += finalWord + " " + strconv.Itoa(delta)
		} else {
			result += "tkhapjiabb" + " " + strconv.Itoa(0)
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
		!strings.Contains(blackWord, string(currentString[l:r])) &&
		l+r != len(currentString) {
		return string(currentString[l:r]), true
	}

	return "", false
}

func countOccurrences(words []string, substr string) int {
	occurrences := 0

	for i := 0; i < len(words); i++ {
		if strings.Contains(words[i], substr) {
			occurrences++
		}
	}

	return occurrences
}
