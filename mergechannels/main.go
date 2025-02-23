package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	channels := make([]chan int64, 10)

	for i := range channels {
		channels[i] = make(chan int64)
	}

	for i := range channels {
		go func(i int) {
			channels[i] <- int64(i)
			close(channels[i])
		}(i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	for v := range merge(ctx, channels...) {
		fmt.Println(v)
	}
}

// FanIn pattern
func merge(ctx context.Context, channels ...chan int64) <-chan int64 {
	result := make(chan int64, 1)

	merger := func(wg *sync.WaitGroup, ch chan int64) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
				return
			case value, ok := <-ch:
				if !ok {
					return
				}

				result <- value
				wg.Done()
			}
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		go merger(wg, channel)
	}

	go func(res chan int64) {
		wg.Wait()
		close(res)
	}(result)

	return result
}
