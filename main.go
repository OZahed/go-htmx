package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	NoColor     = "\u001b[0m"
	RedColor    = "\u001b[41m"
	GreenColor  = "\u001b[42m"
	YellowColor = "\u001b[43m"
	BlueColor   = "\u001b[44m"
	CyanColor   = "\u001b[46m"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags | log.Lmicroseconds)
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("GET /badreq", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	})

	mux.HandleFunc("GET /internalerr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: TimeIt(PanicHandler(mux)),
	}

	log.Println("Server is ready and Litens on port:3000, you can open http://localhost:3000/")
	log.Fatal(server.ListenAndServe())
}

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

func colorizeDuration(d time.Duration) string {
	switch {
	case d < 1*time.Millisecond:
		return GreenColor + " " + d.String() + " " + NoColor
	case d < 3*time.Millisecond:
		return YellowColor + " " + d.String() + " " + NoColor
	default:
		return RedColor + " " + d.String() + " " + NoColor
	}
}

func colorizeStatus(status int) string {
	switch {
	case status > 199 && status < 300:
		return GreenColor + " " + strconv.Itoa(status) + " " + NoColor
	case status > 299 && status < 400:
		return BlueColor + " " + strconv.Itoa(status) + " " + NoColor
	case status > 399 && status < 500:
		return YellowColor + " " + strconv.Itoa(status) + " " + NoColor
	default:
		return RedColor + " " + strconv.Itoa(status) + " " + NoColor
	}
}
