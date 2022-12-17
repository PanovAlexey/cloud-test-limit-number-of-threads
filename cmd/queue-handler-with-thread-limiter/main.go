package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const MAX_THREADS = 2
const MAX_QUEUE_LEN = 4
const MOCK_WORKER_TIME = 15
const SHOUTDOWN_CONTEXT_TIMEOUT = 60

type Task struct {
}

func (t Task) Do() {
	time.Sleep(time.Second * MOCK_WORKER_TIME)

	log.Println("Task handled", time.Now().String())
}

func PostTaskToQueue(queue chan Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			queue <- Task{}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("task added"))
			return
		}

		http.Error(w, "Method is not allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func queueHandler(queue chan Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range queue {
		task.Do()
	}

	log.Println("handler goroutine is closed")
}

func main() {
	queue := make(chan Task, MAX_QUEUE_LEN)
	wg := &sync.WaitGroup{}

	for i := 0; i < MAX_THREADS; i++ {
		wg.Add(1)
		go queueHandler(queue, wg)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", PostTaskToQueue(queue))

	srv := &http.Server{Addr: ":8080"}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("http server error. ", err)
		}
	}()

	log.Println("signal for server shutdown detected: ", <-sigs)

	shoutDownCtx := context.Background()
	shoutDownCtx, cancel := context.WithTimeout(shoutDownCtx, time.Second*SHOUTDOWN_CONTEXT_TIMEOUT)
	defer cancel()

	if err := srv.Shutdown(shoutDownCtx); err != nil {
		log.Println("server shoutdowning error. ", err)
	}

	close(queue)
	wg.Wait()

	log.Println("server has been shutdown")
}
