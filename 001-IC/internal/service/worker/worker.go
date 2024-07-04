package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

// worker, tokenvalidator, пакет работы с файлами
// Worker executes a task at regular intervals
type Worker struct {
	task     func()
	interval time.Duration
}

// NewWorker creates a new instance of Worker
func NewWorker(task func(), interval time.Duration) *Worker {
	return &Worker{
		task:     task,
		interval: interval,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context done, executing task before exit")
			w.task()
			log.Println("Stopping flush to file:", ctx.Err())
			return
		case <-ticker.C:
			log.Println("Ticker ticked, executing task")
			w.task()
		}
	}
}
