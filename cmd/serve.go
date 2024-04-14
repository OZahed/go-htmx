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
	"github.com/OZahed/go-htmx/internal/handlers"
	"github.com/OZahed/go-htmx/internal/handlers/middleware"
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
	tmp := tmpl.LoadTemplates(cfg.LayoutsRootDir)
	lg := logger.NewLogger().With("name", cfg.DebuggerBaseName)
	// make handlers
	layoutHandlers := handlers.NewLayoutHandler(tmp, cfg.AppName, cfg.LayoutRootTmpName, lg)
	// add health check route
	healthHandler := handlers.NewHealthHandler()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(cfg.StaticFilesDir))
	mux.Handle("GET /public/", http.StripPrefix(cfg.StaticRoutesPrefix, fs))

	handlers.SetHTMLRoutes(mux, layoutHandlers)
	handlers.SetHandlerRoutes(mux, healthHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: middleware.TimeIt(middleware.PanicHandler(mux)),
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
