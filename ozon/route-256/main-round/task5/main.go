package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func prettify(input interface{}) interface{} {
	switch data := input.(type) {
	case map[string]interface{}:
		for key, value := range data {
			if cleaned := prettify(value); cleaned == nil {
				delete(data, key)
			} else {
				data[key] = cleaned
			}
		}

		if len(data) == 0 {
			return nil
		}

		return data
	case []interface{}:
		cleanedSlice := make([]interface{}, 0)

		for _, item := range data {
			if cleaned := prettify(item); cleaned != nil {
				cleanedSlice = append(cleanedSlice, cleaned)
			}
		}

		if len(cleanedSlice) == 0 {
			return nil
		}

		return cleanedSlice
	default:
		return data
	}
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	in.Scan()
	numTests, _ := strconv.Atoi(in.Text())

	results := make([]interface{}, numTests)

	for i := 0; i < numTests; i++ {
		in.Scan()
		numLines, _ := strconv.Atoi(in.Text())

		var jsonBuilder strings.Builder
		for j := 0; j < numLines; j++ {
			in.Scan()
			jsonBuilder.WriteString(in.Text())
		}

		var jsonData interface{}
		if err := json.Unmarshal([]byte(jsonBuilder.String()), &jsonData); err != nil {
			fmt.Println(err)
		}

		results[i] = prettify(jsonData)
	}

	finalJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintln(out, string(finalJSON))
	out.Flush()
}
