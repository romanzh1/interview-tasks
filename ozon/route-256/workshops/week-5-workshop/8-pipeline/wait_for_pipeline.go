package main

import "sync"

func waitForPipeline(errs ...<-chan error) error {
	errc := merge(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func merge(cs ...<-chan error) <-chan error {
	out := make(chan error)
	wg := sync.WaitGroup{}

	output := func(c <-chan error) {
		defer wg.Done()

		for m := range c {
			out <- m
		}
	}

	for _, c := range cs {
		wg.Add(1)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
