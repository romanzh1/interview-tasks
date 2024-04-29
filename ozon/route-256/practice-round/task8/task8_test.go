package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateString(t *testing.T) {
	for i := 1; i <= 20; i++ {
		inputFile := strconv.Itoa(i) + ""
		outputFile := strconv.Itoa(i) + ".a"

		input, err := os.Open(inputFile)
		defer input.Close()
		require.NoError(t, err, "Error opening input file: %s", inputFile)

		expectedOutput, err := os.ReadFile(outputFile)
		require.NoError(t, err, "Error opening output file: %s", outputFile)

		scanner := bufio.NewScanner(input)
		result := calculateString(scanner)

		// Разбиваем результаты на строки
		expectedLines := strings.Split(strings.TrimSpace(string(expectedOutput)), "\n")
		resultLines := strings.Split(strings.TrimSpace(result), "\n")

		// Проверяем, что длины массивов одинаковы
		require.Equal(t, len(expectedLines), len(resultLines), "Number of lines mismatch in test case %d", i)

		// Сравниваем построчно
		for index := range expectedLines {
			if expectedLines[index] != resultLines[index] {
				// Если строки не совпадают, сравниваем числа в конце строк
				expectedParts := strings.Split(expectedLines[index], " ")
				resultParts := strings.Split(resultLines[index], " ")

				// Сравниваем числа
				expectedNumber := expectedParts[len(expectedParts)-1]
				resultNumber := resultParts[len(resultParts)-1]
				if expectedNumber != resultNumber {
					testCaseContent := readTestCase(t, inputFile, index+1)
					require.Equal(t, expectedLines[index], resultLines[index],
						"Mismatch in line %d of test case %d:\n%s", index+1, i, testCaseContent)
				}
			}
		}
	}
}

func readTestCase(t *testing.T, filePath string, testCaseNumber int) string {
	t.Helper()

	file, err := os.Open(filePath)
	require.NoError(t, err, "Error opening input file: %s", filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	numberOfTest := 0

	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Fields(line)
		_, err := strconv.Atoi(args[0])
		if err != nil {
			continue
		}

		if numberOfTest == testCaseNumber {
			content += line + "\n"
			for scanner.Scan() {
				line := scanner.Text()
				args := strings.Fields(line)
				_, err := strconv.Atoi(args[0])
				if err == nil {
					return content
				}

				content += line + "\n"
			}
		}

		numberOfTest++
	}

	return content
}
