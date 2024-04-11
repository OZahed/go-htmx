package main

import (
	"net/http"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}
