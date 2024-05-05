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
	t, _ := strconv.Atoi(in.Text())
	results := make([]float64, t)

	for i := 0; i < t; i++ {
		bankRates := make([][6]float64, 3)

		for j := 0; j < 3; j++ {
			for k := 0; k < 6; k++ {
				in.Scan()
				rates := strings.Split(in.Text(), " ")
				n, _ := strconv.ParseFloat(rates[0], 64)
				m, _ := strconv.ParseFloat(rates[1], 64)
				bankRates[j][k] = m / n
			}
		}

		maxDollars := 0.0
		for a := 0; a < 3; a++ {
			for b := 0; b < 3; b++ {
				for c := 0; c < 3; c++ {
					if a != b && b != c && a != c {
						directUSD := bankRates[a][0]
						rubToEurToUSD := bankRates[a][1] * bankRates[b][5]
						rubToUSDToRubToUSD := bankRates[a][0] * bankRates[b][2] * bankRates[c][0]
						rubToUSDToEURToUSD := bankRates[a][0] * bankRates[b][3] * bankRates[c][5]
						rubToEURToRubToUSD := bankRates[a][1] * bankRates[b][4] * bankRates[c][0]

						maxDollars = max(maxDollars, directUSD)
						maxDollars = max(maxDollars, rubToEurToUSD)
						maxDollars = max(maxDollars, rubToUSDToRubToUSD)
						maxDollars = max(maxDollars, rubToUSDToEURToUSD)
						maxDollars = max(maxDollars, rubToEURToRubToUSD)
					}
				}
			}
		}

		results[i] = maxDollars
	}

	fmt.Fprintln(out, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(results)), "\n"), "[]"))
	out.Flush()
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}

	return b
}
