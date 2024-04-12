package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OZahed/go-htmx/internal/log"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (s *StatusRecorder) WriteHeader(statuCode int) {
	s.Status = statuCode
	s.ResponseWriter.WriteHeader(statuCode)
}

func TimeIt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		next.ServeHTTP(recorder, r)
		fmt.Printf("%s |  %s  |  %-10s  |  %s\n", t.Format(time.RFC3339),
			log.ColorizeStatus(recorder.Status), log.ColorizeDuration(time.Since(t)), r.URL.Path,
		)
	})
}
