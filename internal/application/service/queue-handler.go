package service

import (
	"cloud-test-limit-the-number-of-threads/internal/domain"
	"log"
	"sync"
)

const MAX_THREADS = 2

type queueHandler struct {
}

func GetQueueHandler() queueHandler {
	return queueHandler{}
}

func (h queueHandler) Run(queue chan domain.Task, wg *sync.WaitGroup) {
	for i := 0; i < MAX_THREADS; i++ {
		wg.Add(1)
		go h.handle(queue, wg)
	}
}

func (h queueHandler) handle(queue chan domain.Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range queue {
		task.Do()
	}

	log.Println("handler goroutine is closed")
}
