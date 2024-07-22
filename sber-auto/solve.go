package main

import (
	"context"
	"fmt"
	"io"
	"sync"
)

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

// 5. Не завершаем го рутину, если на записываем в канал ошибку. Можем потёчь по памяти
// Не закрываем канал после записи.
// Правильнее сделать например так:

func (s *Service) ProcessData2(timeoutCtx context.Context, r io.Reader) error {
	errCh := make(chan error, 1) // Используем буферизированный канал

	go func() {
		errCh <- s.processDataInternal(r)
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

// 6. Функция isCallAllowed ограничивает количество вызовов в течение текущего и предыдущего 30-секундного интервала.
// Она возвращает true, если вызов разрешен, и false в противном случае
