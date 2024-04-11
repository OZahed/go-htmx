package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags | log.Lmicroseconds)
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /about", aboutHandler)
	mux.HandleFunc("GET /blog", blogHandler)
	mux.HandleFunc("GET /tags", tagsHandler)

	server := http.Server{
		Addr:    ":3000",
		Handler: TimeIt(PanicHandler(mux)),
	}

	log.Println("Server is ready and Litens on port:3000, you can open http://localhost:3000/")
	log.Fatal(server.ListenAndServe())
}
