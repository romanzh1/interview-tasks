package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// https://code.yandex-team.ru/041682aa-75a5-4789-bffa-46d8d40b56b4

// 1. Что выведет код и почему?

func setLinkHome(link *string) {
	*link = "http://home"
}

func main1() {
	link := "http://other"

	setLinkHome(&link)

	fmt.Println(link)
}

// 2. Будет ли напечатан “ok” ?

func main2() {
	defer func() {
		recover()
	}()

	panic("test panic")

	fmt.Println("ok")
}

// 3.

// Функция должна напечатать:
// one
// two
// three
// (в любом порядке и в конце обязательно)
// done!

// Но это не так, исправь код

func printText3(data []string) {
	for _, v := range data {
		go func() {
			fmt.Println(v)
		}()

	}

	fmt.Println("done!")

}

func main3() {
	data := []string{"one", "two", "three"}

	printText(data)
}

// 4.
// Мы пытаемся подсчитать количество выполненных параллельно операций,
// что может пойти не так?

var callCounter uint

func main4() {
	wg := sync.WaitGroup{}

	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			// Ходим в базу, делаем долгую работу
			time.Sleep(time.Second)
			// Увеличиваем счетчик
			callCounter++
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Call counter value = ", callCounter)
}

// 5.
// Есть функция processDataInternal, которая может выполняться неопределенно долго.
// Чтобы контролировать процесс, мы добавили таймаут выполнения ф-ии через context.
// Какие недостатки кода ниже?

type Service struct {
}

func (s *Service) processDataInternal(r io.Reader) error {
	return nil
}

func (s *Service) ProcessData(timeoutCtx context.Context, r io.Reader) error {
	errCh := make(chan error)

	go func() {
		errCh <- s.processDataInternal(r)
	}()

	select {
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

// 6.
// Опиши, что делает функция isCallAllowed?
// TODO разобрать что же она всё таки делает

var callCount = make(map[uint]uint)

var locker = &sync.Mutex{}

func isCallAllowed(allowedCount uint) bool {
	if allowedCount == 0 {
		return true
	}

	locker.Lock()
	defer locker.Unlock()

	curTimeIndex := uint(time.Now().Unix() / 30)

	prevIndexVal, _ := callCount[curTimeIndex-1]

	if prevIndexVal >= allowedCount {
		return false
	}

	curIndexVal, ok := uint(0), false
	if curIndexVal, ok = callCount[curTimeIndex]; !ok {
		callCount[curTimeIndex] = 1

		return true
	}

	if (curIndexVal + prevIndexVal) >= allowedCount {
		return false
	}

	callCount[curTimeIndex]++

	return true
}

func main() {
	fmt.Printf("%v\n", isCallAllowed(3)) // true
	fmt.Printf("%v\n", isCallAllowed(3)) // true
	fmt.Printf("%v\n", isCallAllowed(3)) // true

	// time.Sleep(time.Second*30)

	fmt.Printf("%v\n", isCallAllowed(3)) // false
	fmt.Printf("%v\n", isCallAllowed(3)) // false

}

// 9. Нарисовать архитектуру приложения, что будет происходить на бэкенде при нажатии кнопки пользователем
// В задаче есть 2 внешние сущности, эндпоинт, возвращающий пользователей и банк
