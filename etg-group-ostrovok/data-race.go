package main

package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// FIXME:
func main() {    
    rates, err := readConversionRates()
    if err != nil {
        log.Fatalf("read intial conversion rates values: %s", err)
    }

    // background task to update conversion rates
    go func() {
        for {
            time.Sleep(time.Minute)
            rates, err = readConversionRates()
            if err{
                log.Printf("ERR: update converson rates: %s", err)
                continue                
            }
        }
    }()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // simplified, for short
        from := "RUB"
        val := 140.0

        rate, _ := rates[from]

        convVal := val / rate

        fmt.Fprint(w, convVal)
    })

    if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}

// solve
func main() {
	// add a mutex to read from the map
    mu := sync.Mutex{}
    
    mu.Lock()
    rates, err := readConversionRates()
    if err != nil {
        log.Fatalf("read intial conversion rates values: %s", err)
    }
    mu.Unlock()

    // background task to update conversion rates
    go func() {
        for {
            time.Sleep(time.Minute)
            tmp, err := readConversionRates() // long operation
            if err{
                log.Printf("ERR: update converson rates: %s", err)
                continue                
            }
            mu.Lock()
            rates = tmp
            mu.Unlock()
        }
    }()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // simplified, for short
        from := "RUB"
        val := 140.0

        mu.Lock() ////
        rate, _ := rates[from]
        mu.Unlock()

        convVal := val / rate

        fmt.Fprint(w, convVal)
    })

    if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}

// readConversionRates reads rates from a file or an external service (relatively long-running function).
func readConversionRates() (map[string]float64, error) { 
    time.Sleep(100*time.Milliseconds) 
    return map[string]float64{
        "USD": 1.0,
        "RUB": 70.0,
    }, nil
}