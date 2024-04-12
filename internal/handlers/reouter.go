package handlers

import "net/http"

func SetHTMLRoutes(router *http.ServeMux, handler *LayoutHandlers) {
	router.HandleFunc("GET /", handler.IndexHandler)
	router.HandleFunc("GET /about", handler.AboutHandler)
	router.HandleFunc("GET /blog", handler.BlogHandler)
	router.HandleFunc("GET /tags", handler.TagsHandler)
}
