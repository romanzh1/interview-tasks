package main

import (
	"context"
	"fmt"
)

// []string -> parsing -> if n < 100 -> reduce

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errcList []<-chan error

	input := []string{"1", "12", "25", "1xx50"}

	//step 1
	linesCh, errc, err := source(ctx, input...)
	if err != nil {
		panic(err)
	}
	errcList = append(errcList, errc)

	//step 2
	numberCh, errc := parse(ctx, linesCh)
	errcList = append(errcList, errc)

	//step3
	resc, errc := sink(ctx, numberCh)
	errcList = append(errcList, errc)

	fmt.Println("started")

	if err = waitForPipeline(errcList...); err != nil {
		fmt.Println("ERR: ", err)
		cancel()
		return
	}

	res := <-resc
	fmt.Printf("result is: %d\n", res)
	fmt.Println("close app")
}
