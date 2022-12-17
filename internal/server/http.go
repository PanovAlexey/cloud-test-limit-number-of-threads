package server

import (
	"cloud-test-limit-the-number-of-threads/internal/domain"
	"cloud-test-limit-the-number-of-threads/internal/endpoint"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const SHOUTDOWN_CONTEXT_TIMEOUT = 60

type mainHttpServer struct {
	httpServer *http.Server
}

func RunHttpServer(queue chan domain.Task) {
	mux := http.NewServeMux()
	mux.Handle("/add", endpoint.PostTaskToQueue(queue))
	server := mainHttpServer{httpServer: &http.Server{Addr: ":8080", Handler: mux}}

	go func() {
		log.Println("server starting...")

		if err := server.httpServer.ListenAndServe(); err != nil {
			log.Println("http server error. ", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("signal for server shutdown detected: ", <-sigs)

	shoutDownCtx := context.Background()
	shoutDownCtx, cancel := context.WithTimeout(shoutDownCtx, time.Second*SHOUTDOWN_CONTEXT_TIMEOUT)
	defer cancel()

	if err := server.httpServer.Shutdown(shoutDownCtx); err != nil {
		log.Println("server shoutdowning error. ", err)
	}

	log.Println("server has been shutdown")
}
