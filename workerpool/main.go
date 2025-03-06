package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID       int64
	Filename string
}

type WorkerPool struct {
	numWorkers int
	numTasks   int
	taskCh     chan Task
	resCh      chan string
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, numTasks int) *WorkerPool {
	wp := &WorkerPool{
		numWorkers: numWorkers,
		numTasks:   numTasks,
		taskCh:     make(chan Task, numTasks),
		resCh:      make(chan string, numTasks),
	}

	wp.wg = sync.WaitGroup{}
	wp.wg.Add(numWorkers)

	return wp
}

func (wp *WorkerPool) runWorker(workerID int) {
	for task := range wp.taskCh {
		fmt.Printf("Worker %d started processing task %d\n", workerID, task.ID)
		wp.resCh <- ProcessFile(task)
		fmt.Printf("Worker %d finished processing task %d\n", workerID, task.ID)
	}
}

func (wp *WorkerPool) RunWorkers() {
	for i := 0; i < wp.numWorkers; i++ {
		go func(i int) {
			defer wp.wg.Done()
			wp.runWorker(i)
		}(i)
	}
}

func (wp *WorkerPool) WaitAllWorkers() {
	go func() {
		wp.wg.Wait()
		close(wp.resCh)
	}()
}

func (wp *WorkerPool) CreateTasks() {
	go func() {
		for i := 0; i < wp.numTasks; i++ {
			wp.taskCh <- Task{
				ID:       int64(i),
				Filename: fmt.Sprintf("File_%d", i+1),
			}
		}
		close(wp.taskCh)
	}()
}

func ProcessFile(task Task) string {
	time.Sleep(time.Second)
	return fmt.Sprintf("%s processed (Task ID: %d)",
		task.Filename, task.ID)
}

func main() {
	const (
		numWorkers = 3
		numTasks   = 10
	)

	wp := NewWorkerPool(numWorkers, numTasks)

	wp.RunWorkers()
	wp.CreateTasks()
	wp.WaitAllWorkers()

	for res := range wp.resCh {
		fmt.Println(res)
	}

	fmt.Println("All tasks processed")
}
