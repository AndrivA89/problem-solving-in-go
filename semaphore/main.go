package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	wg    sync.WaitGroup
	slots chan struct{}
}

func NewSemaphore(amount int) *Semaphore {
	return &Semaphore{
		slots: make(chan struct{}, amount),
	}
}

func (s *Semaphore) Get() {
	s.wg.Add(1)
	s.slots <- struct{}{}
}

func (s *Semaphore) Free() {
	s.wg.Done()
	<-s.slots
}

func (s *Semaphore) Run(f func()) {
	go func() {
		s.Get()
		f()
		s.Free()
	}()
}

func (s *Semaphore) WaitDone() {
	s.wg.Wait()
	close(s.slots)
}

func DownloadFile(filepath string) {
	fmt.Printf("Downloading %s...\n", filepath)
	time.Sleep(time.Second)
	fmt.Printf("Downloaded %s\n", filepath)
}

func main() {
	semaphore := NewSemaphore(3)

	for i := 0; i < 10; i++ {
		file := fmt.Sprintf("File_%d.txt", i+1)
		semaphore.Run(func() { DownloadFile(file) })
	}

	semaphore.WaitDone()
}
