package main

import (
	"net/http"
	"time"
)

const MAX_THREADS = 8
const MAX_QUEUE_LEN = 1_000

type Task struct {
}

func (t Task) Do() {
	time.Sleep(time.Second * 10)
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

func queueHandler(queue chan Task) {
	for true {
		task := <-queue
		task.Do()
	}
}

func main() {
	queue := make(chan Task, MAX_QUEUE_LEN)

	for i := 0; i < MAX_THREADS; i++ {
		go queueHandler(queue)
	}

	http.HandleFunc("/", PostTaskToQueue(queue))
	http.ListenAndServe(":8080", nil)
}
