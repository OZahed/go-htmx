package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Combine(mux http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return mux
	}

	for _, m := range middlewares {
		mux = m(mux)
	}

	return mux
}
