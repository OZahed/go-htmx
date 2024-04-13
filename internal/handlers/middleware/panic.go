package middleware

import (
	"log"
	"net/http"
)

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic did happen with error: ", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
