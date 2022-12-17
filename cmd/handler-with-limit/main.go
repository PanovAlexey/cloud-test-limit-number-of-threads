package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MAX_THREADS = 8
const MAX_QUEUE_LEN = 1_000

type Task struct {
}

func (t Task) Do() {
	time.Sleep(time.Second * 10)

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", PostTaskToQueue(queue))

	srv := &http.Server{Addr: ":8080"}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("Http server error. ", err)
		}
	}()

	log.Println("Signal for server shutdown detected: ", <-sigs)

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Println("server shoutdowning error")
	}

	close(queue)

	log.Println("Server has been shutdown")
}
