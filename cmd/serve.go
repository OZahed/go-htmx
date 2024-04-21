package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OZahed/go-htmx/internal/config"
	"github.com/OZahed/go-htmx/internal/handler"
	"github.com/OZahed/go-htmx/internal/handler/middleware"
	"github.com/OZahed/go-htmx/internal/handler/tmpl"
	"github.com/OZahed/go-htmx/internal/logger"
)

var (
	_ Command = (*ServeCmd)(nil)
)

type ServeCmd struct{}

// Help implements Command.
func (ServeCmd) Help() HelpInfo {
	return HelpInfo{
		SubCmdName: "serve",
		ShortDesc:  "runs server for the application",
		LongDesc:   "",
		Usage:      fmt.Sprintf("%s serve", APP_NAME),
	}
}

// Name implements Command.
func (ServeCmd) Name() string {
	return "serve"
}

func (s ServeCmd) Execute(args []string) {
	ver := VersionCmd{}
	ver.Execute(nil)

	cfg, err := config.NewAppConfig(APP_NAME)
	if err != nil {
		log.Fatal(err)
	}

	layoutTemp := tmpl.LoadTemplates(cfg.Layout.TempDir)
	partialTemp := tmpl.LoadTemplates(cfg.Static.PartialsDir)

	lg := logger.NewLogger()
	// make handlers
	layoutHandler := handler.NewLayout(layoutTemp, cfg.AppName, cfg.Layout.TempRootName, lg)
	partialHandler := handler.NewPartials(partialTemp, lg.With("name", "partials"))
	// add health check route
	healthHandler := handler.NewHealthHandler()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(cfg.Static.FilesDir))
	mux.Handle("GET /public/", middleware.GZip(1)(
		http.StripPrefix(cfg.Static.RoutesPrefix, fs),
	))

	handler.SetHTMLRoutes(mux, layoutHandler)
	handler.SetHandlerRoutes(mux, healthHandler)
	handler.SetPartialRoute(mux, partialHandler)

	// Creating a sub route
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello from apis requested for id: " + r.PathValue("id")))
	})
	apiMux.HandleFunc("POST /todos/{id}/done", func(w http.ResponseWriter, r *http.Request) {})
	apiMux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {})

	apiMuxMiddlewares := []middleware.Middleware{}
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", middleware.Combine(apiMux, apiMuxMiddlewares...)))

	// middlewares apply in reverse order
	middlewares := []middleware.Middleware{
		middleware.PanicHandler,
		middleware.LogIt,
		middleware.GZip(3),
	}

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           middleware.Combine(mux, middlewares...),
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	ctx, cnl := signal.NotifyContext(context.Background(),
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cnl()

	go func() {
		var err error
		if cfg.Server.CertFile != "" {
			lg.Debug(fmt.Sprintf("Server is ready -> https://localhost:%d/", cfg.Server.Port))
			err = server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile)
		} else {
			lg.Debug(fmt.Sprintf("Server is ready -> http://localhost:%d/", cfg.Server.Port))
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error("failed to run server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	lg.Warn("got signal for shutdown, server is shutting down ...")

	shutCtx, shutCnl := context.WithDeadline(context.Background(), time.Now().Add(cfg.Server.ShutdownDuration))
	defer shutCnl()

	if err := server.Shutdown(shutCtx); err != nil {
		lg.ErrorContext(shutCtx, "server shutdown got error", "error", err)
	}

	lg.Warn("server is down")
}
