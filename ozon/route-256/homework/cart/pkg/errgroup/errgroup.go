package errgroup

import (
	"sync"
)

type Group struct {
	wg      sync.WaitGroup
	err     error
	errOnce sync.Once
	sem     chan struct{}
}

func NewGroup(limit int) *Group {
	return &Group{
		sem: make(chan struct{}, limit),
	}
}

func (g *Group) done() {
	<-g.sem
	g.wg.Done()
}

func (g *Group) Go(f func() error) {
	g.sem <- struct{}{}
	g.wg.Add(1)
	go func() {
		defer g.done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()

	return g.err
}
