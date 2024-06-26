package handler

import "net/http"

func SetHTMLRoutes(router *http.ServeMux, handler *Layout) {
	router.HandleFunc("/", handler.IndexHandler)
	router.HandleFunc("GET /about", handler.AboutHandler)
	router.HandleFunc("GET /blog", handler.BlogHandler)
	router.HandleFunc("GET /tags", handler.TagsHandler)
	router.HandleFunc("GET /go/greetings/{name}/{id}", handler.Greet)
}

func SetHandlerRoutes(router *http.ServeMux, handler *HealthHandler) {
	router.HandleFunc("GET /health", handler.Health)
}

func SetPartialRoute(router *http.ServeMux, handler *Partials) {
	router.HandleFunc("GET /counter/{count}", handler.Counter)
}

func SetTemplRoutes(router *http.ServeMux, handler *TemplHandler) {
	router.HandleFunc("GET /templ/greetings/{name}/{age}", handler.Greet)
}
