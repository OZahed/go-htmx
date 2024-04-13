package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/OZahed/go-htmx/internal/logger"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status   int
	ByteSize int
}

func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (s *StatusRecorder) Write(b []byte) (n int, err error) {
	s.ByteSize = s.ByteSize + len(b)
	return s.ResponseWriter.Write(b)
}

func TimeIt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		next.ServeHTTP(recorder, r)
		log.Printf("%s | %-15s | %-10s | %s",
			logger.ColorizeStatus(recorder.Status),
			logger.ColorizeDuration(time.Since(t)),
			logger.HumanReadableBytes(recorder.ByteSize),
			r.URL.Path,
		)
	})
}
