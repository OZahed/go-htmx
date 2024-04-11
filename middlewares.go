package main

import (
	"log"
	"net/http"
	"time"
)

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("paniced with error: ", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func TimeIt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		next.ServeHTTP(recorder, r)
		log.Printf("|  %s  |  %s  |  %s\n", colorizeStatus(recorder.Status), colorizeDuration(time.Since(t)), r.URL.Path)
	})
}
