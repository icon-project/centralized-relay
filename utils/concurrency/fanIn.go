package concurrency

import "sync"

func FanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	multiplex := func(c <-chan interface{}, wg *sync.WaitGroup, mStream chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case mStream <- i:
			}
		}
	}

	multiplexedStream := make(chan interface{})
	wg := &sync.WaitGroup{}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c, wg, multiplexedStream)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
