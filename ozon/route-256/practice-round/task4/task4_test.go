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
	for i := 1; i <= 50; i++ {
		inputFile := "test_files/" + strconv.Itoa(i)
		outputFile := "test_files/" + strconv.Itoa(i) + ".a"

		input, err := os.Open(inputFile)
		defer input.Close()
		require.NoError(t, err, "Error opening input file: %s", inputFile)

		expectedOutput, err := os.ReadFile(outputFile)
		require.NoError(t, err, "Error opening output file: %s", outputFile)

		reader := bufio.NewReader(input)
		result := calculateString(reader)

		require.Equal(t, strings.Split(string(expectedOutput), " \n"), strings.Split(result, "\n"),
			"Test case %d", i)
	}
}
