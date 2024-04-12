package cmd

import (
	"net/http"
	"os"

	"github.com/OZahed/go-htmx/internal/handlers"
	"github.com/OZahed/go-htmx/internal/handlers/middleware"
	"github.com/OZahed/go-htmx/internal/log"
)

func ExucuteServe() {
	// read configs
	// make templates
	tmp := handlers.LoadTemplates("./templates")
	lg := log.NewLogger().With("name", "main")
	// make handlers
	layoutHandlers := handlers.NewLayoutHanler(tmp, "OZahed", "Layout", lg)
	// add health check route
	// add routes the serveMux
	// listen and serve
	// add graceful shutdown

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	handlers.SetHTMLRoutes(mux, layoutHandlers)

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.TimeIt(middleware.PanicHandler(mux)),
	}

	lg.Debug("Server is ready and Litens on port:3000, you can open http://localhost:3000/")
	err := server.ListenAndServe()
	if err != nil {
		lg.Error("failed to run server", err)
		os.Exit(1)
	}
}