package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const ()

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZip(level int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			// it is better to use a writer pool instead of just create one and defer it's close on every connection
			gz := gzip.NewWriter(w)
			defer func() {
				if err := gz.Close(); err != nil {
					panic(err)
				}
			}()

			gzResponseWriter := gzipResponseWriter{Writer: gz, ResponseWriter: w}

			next.ServeHTTP(gzResponseWriter, r)
		})
	}
}
