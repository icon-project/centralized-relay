package concurrency

import "sync"

/*
FanIn multiplexes multiple input channels into a single output channel. It reads from each input channel
concurrently until all channels are closed, then closes the output channel. The function takes a 'done'
channel to signal cancellation, and variadic 'channels' representing the input channels to multiplex.
It returns a single output channel where values from all input channels are sent.
*/
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

/*
Take reads values from the 'valueStream' input channel and sends them to the output channel 'takeStream'.
It reads until 'noOfItemsToTake' items are taken from 'valueStream' or a cancellation signal is received
on the 'done' channel. The function returns a channel 'takeStream' where the taken values are sent.
*/
func Take(done <-chan interface{}, valueStream <-chan interface{}, noOfItemsToTake int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < noOfItemsToTake; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}
