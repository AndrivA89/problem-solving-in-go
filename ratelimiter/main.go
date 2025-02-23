package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(chan struct{}, rps),
	}

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rps))
		defer ticker.Stop()

		for range ticker.C {
			select {
			case rl.requests <- struct{}{}:
			default: // Если канал заполнен, пропускаем (чтобы не забивать буфер).
			}
		}
	}()

	return rl
}

func (r *RateLimiter) someRequest(url string) error {
	select {
	case <-r.requests:
		response, err := http.Get(url)
		if err != nil {
			return err
		}

		if err = response.Body.Close(); err != nil {
			return err
		}

		return nil
	case <-time.After(time.Second): // Ожидание 1 сек перед тем как отвалиться.
		return errors.New("rate limit exceeded (timeout)")
	}
}

func main() {
	// Разрешаем делать 5 запросов в секунду.
	rl := NewRateLimiter(5)
	// Подождём, пока заполнится канал токенами.
	time.Sleep(2 * time.Second)

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

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			if err := rl.someRequest(u); err != nil {
				fmt.Printf("Request to %s failed: %s\n", u, err)
			} else {
				fmt.Printf("Request to %s OK\n", u)
			}
		}(url)
	}

	wg.Wait()
}
