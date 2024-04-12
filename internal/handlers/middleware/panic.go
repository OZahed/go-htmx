package middleware

import (
	"log"
	"net/http"
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
