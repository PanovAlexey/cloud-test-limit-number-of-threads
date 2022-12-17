package endpoint

import (
	"cloud-test-limit-the-number-of-threads/internal/domain"
	"net/http"
)

func PostTaskToQueue(queue chan domain.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			queue <- domain.Task{}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("task added"))
			return
		}

		http.Error(w, "Method is not allowed!", http.StatusMethodNotAllowed)
		return
	}
}
