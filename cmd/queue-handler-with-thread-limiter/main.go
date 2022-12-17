package main

import (
	"cloud-test-limit-the-number-of-threads/internal/application/service"
	"cloud-test-limit-the-number-of-threads/internal/domain"
	server "cloud-test-limit-the-number-of-threads/internal/server"
	"log"
	"sync"
)

const MAX_QUEUE_LEN = 4

func main() {
	queue := make(chan domain.Task, MAX_QUEUE_LEN)
	wg := &sync.WaitGroup{}

	queueHandler := service.GetQueueHandler()
	queueHandler.Run(queue, wg)

	server.RunHttpServer(queue)

	close(queue)
	wg.Wait()

	log.Println("Application was closed.")
}
