package main

// task to implement the function
func merge(chs ...chan int) chan int {

}

func main() {
	ch1 := startProducerA()
	ch2 := startProducerB()

	for el := range merge(ch1, ch2) {
		println(el)
	}
}
