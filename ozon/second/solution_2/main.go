package main

import (
	"fmt"
	"math/rand"
	"sync"
)

//////////////////

// Поочередно выполнит http запросы по предложенному списку ссылок

// в случае получения http-кода ответа на запрос "200 OK" печатаем на экране "адрес url - ok"
// в случае получения http-кода ответа на запрос отличного от "200 OK" либо в случае ошибки печатаем на экране "адрес url - not ok"

// Модифицируйте программу таким образом, чтобы использовались каналы для коммуникации основного потока с горутинами. Пример:

// Запросы по списку выполняются в горутинах.
// Печать результатов на экран происходит в основном потоке

func Get(url string) (resp, error) {
		num := rand.Int() % 10
		if num > 7 {
			return resp{200}, nil
		}
		if num < 4 {
			return resp{404}, nil
		}

		return resp{500}, fmt.Errorf("error fetching url %s", url)
	}

type resp struct {
	code int
}

func main() {
	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://ёёёё",
	}

	wg := sync.WaitGroup{}

	result := make(chan string, len(urls))
	wg.Add(len(urls))
	for i := 0; i < len(urls); i++ {
		go func(i int) {
			defer wg.Done()
			result <- worker(urls[i])
		}(i)
	}

	wg.Wait()
	close(result)

	for l := range result{
		fmt.Println(l)
	}
}

func worker(url string) string {
	res, err := Get(url)
	if err != nil {
		return fmt.Sprintf("url %s is not ok", url)
	}

	if res.code != 200{
		return fmt.Sprintf("url %s is not ok", url)
	}

	return fmt.Sprintf("url %s is ok", url)
}
