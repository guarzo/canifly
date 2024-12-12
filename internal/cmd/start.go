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

	"github.com/guarzo/canifly/internal/server"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var version = "0.0.40"

func Start() error {
	logger := server.SetupLogger()
	logger.Infof("Starting application, version %s", version)

	cfg, err := server.LoadConfig(logger)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	services, err := server.GetServices(logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	r := server.SetupHandlers(cfg.SecretKey, logger, services)
	srv := createServer(r, cfg.Port)

	logger.Infof("Server listening on port %s", cfg.Port)
	return runServer(srv, logger)
}

func createServer(r http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}

func runServer(srv *http.Server, logger interfaces.Logger) error {
	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("Server forced to shutdown")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen: %w", err)
	}

	logger.Info("Server exited cleanly")
	return nil
}
