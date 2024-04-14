package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OZahed/go-htmx/internal/config"
	"github.com/OZahed/go-htmx/internal/handler"
	"github.com/OZahed/go-htmx/internal/handler/middleware"
	"github.com/OZahed/go-htmx/internal/logger"
	"github.com/OZahed/go-htmx/internal/tmpl"
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

	cfg := config.NewAppConfig(APP_NAME)
	layoutTemp := tmpl.LoadTemplates(cfg.LayoutsRootDir)
	partialTemp := tmpl.LoadTemplates(cfg.PartialRootDirs)

	lg := logger.NewLogger()
	// make handlers
	layoutHandler := handler.NewLayout(layoutTemp, cfg.AppName, cfg.LayoutRootTmpName, lg)
	partialHandler := handler.NewPartials(partialTemp, lg.With("name", "partials"))
	// add health check route
	healthHandler := handler.NewHealthHandler()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(cfg.StaticFilesDir))
	mux.Handle("GET /public/", http.StripPrefix(cfg.StaticRoutesPrefix, fs))

	handler.SetHTMLRoutes(mux, layoutHandler)
	handler.SetHandlerRoutes(mux, healthHandler)
	handler.SetPartialRoute(mux, partialHandler)

	// middlewares apply in reverse order
	middlewares := []middleware.Middleware{
		middleware.PanicHandler,
		middleware.LogIt,
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: middleware.Combine(mux, middlewares...),
	}

	ctx, cnl := signal.NotifyContext(context.Background(),
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cnl()

	go func() {
		lg.Debug("Server is ready and Listens on port:3000, you can open http://localhost:3000/")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error("failed to run server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	lg.Warn("got signal for shutdown, server is shutting down ...")

	shutCtx, shutCnl := context.WithDeadline(context.Background(), time.Now().Add(cfg.ShutdownDuration))
	defer shutCnl()

	if err := server.Shutdown(shutCtx); err != nil {
		lg.ErrorContext(shutCtx, "server shutdown got error", "error", err)
	}

	lg.Warn("server is down")
}
