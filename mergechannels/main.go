package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	a := make(chan int64)
	b := make(chan int64)
	c := make(chan int64)

	go func() {
		for _, num := range []int64{1, 2, 3} {
			a <- num
		}
		close(a)
	}()

	go func() {
		for _, num := range []int64{20, 10, 30} {
			b <- num
		}
		close(b)
	}()

	go func() {
		for _, num := range []int64{300, 200, 100} {
			c <- num
		}
		close(c)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	for v := range merge(ctx, a, b, c) {
		fmt.Println(v)
	}
}

// FanIn pattern
func merge(ctx context.Context, channels ...chan int64) <-chan int64 {
	result := make(chan int64, 1)

	merger := func(ch chan int64) {
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
			}
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		go func(ch chan int64) {
			defer wg.Done()
			merger(ch)
		}(channel)
	}

	go func(res chan int64) {
		wg.Wait()
		close(res)
	}(result)

	return result
}
