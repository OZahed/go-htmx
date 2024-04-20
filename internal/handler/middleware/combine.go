package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Combine(mux http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return mux
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		mux = middlewares[i](mux)
	}

	return mux
}
