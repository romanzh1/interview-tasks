// 1. Код выведет http://home. Потому что в функции меняется значение по ссылке

// 2. Не будет

// 3.
func printText(data []string) {

	wg := sync.WaitGroup{}

	wg.Add(len(data))

	for _, v := range data {

		go func(v string) {

			fmt.Println(v)

			wg.Done()

		}(v)

	}

	wg.Wait()

	fmt.Println("done!")

}

// 4. Ответ data race

// 5. Не контролируем канал после выхода из функции, при множественном вызове
// ProcessData приложение может съедать много памяти и даже упасть.