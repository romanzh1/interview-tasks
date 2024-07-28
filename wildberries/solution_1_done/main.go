package main

import (
	"fmt"
	"time"
)

type Agent struct {
	ID      int
	Enabled bool
}

func (a *Agent) Enable() {
	a.Enabled = true
}

type Enabler interface {
	Enable()
}

// 1. Инициализация слайса агентов
// 2. Дополнение слайса сторонними агентами
// 3. Потоковая обработка объектов - активация и отправка на выполнение работы
// 4. Потоковая обработка объектов - сохранение в БД и распечатка резульатов
func main() {
	agents := make([]Agent, 0, 5)
	for i := 0; i < 2; i++ {
		agents = append(agents, Agent{ID: i})
	}

	addThirdPartyAgents(&agents)

	pipe := make(chan Enabler)
	done := make(chan struct{})
	go pipeEnableAndSend(pipe, agents)
	go pipeProcess(pipe, done)

	<-done
	fmt.Println("final")
	close(done)
}

func addThirdPartyAgents(agents *[]Agent) {
	thirdParty := []Agent{
		{ID: 4},
		{ID: 5},
	}
	*agents = append(*agents, thirdParty...)
}

func pipeEnableAndSend(pipe chan Enabler, agents []Agent) {
	for _, a := range agents {
		a := a
		a.Enable()
		pipe <- &a
	}
	close(pipe)
}

func pipeProcess(pipe chan Enabler, done chan struct{}) {
	for a := range pipe {
		dbWrite(a)
	}
	done <- struct{}{}
}

var dbWrite = func(a any) {
	fmt.Println(a)
	time.Sleep(time.Second * 1)
}
