package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guarzo/canifly/internal/server"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var version = "0.0.40"

// Start initializes and runs the application with enhanced logging for startup failure scenarios.
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
	srv, listener, err := createServerWithListener(r, cfg.Port, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	logger.Infof("Server successfully bound on port %s", cfg.Port)
	return runServer(srv, listener, logger)
}

// createServerWithListener attempts to bind to the provided port before creating the HTTP server.
// This helps catch binding errors early.
func createServerWithListener(r http.Handler, port string, logger interfaces.Logger) (*http.Server, net.Listener, error) {
	addr := ":" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.WithError(err).Errorf("Failed to bind to address %s", addr)
		return nil, nil, err
	}
	srv := &http.Server{
		Handler: r,
	}
	return srv, listener, nil
}

// runServer starts the HTTP server and adds enhanced logging to capture startup and shutdown details.
// In DEV_MODE, the startup timeout is increased to give the backend more time to initialize.
func runServer(srv *http.Server, listener net.Listener, logger interfaces.Logger) error {
	startTime := time.Now()
	shutdownCh := make(chan error, 1)

	// Determine startup timeout based on DEV_MODE.
	timeoutDuration := 30 * time.Second
	if os.Getenv("DEV_MODE") == "true" {
		timeoutDuration = 60 * time.Second
		logger.Warn("DEV_MODE is enabled; increasing server startup timeout to 60 seconds")
	}

	// Start the server in a separate goroutine.
	go func() {
		logger.Infof("Server starting to serve requests...")
		err := srv.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("Server encountered an error while serving")
			shutdownCh <- err
		} else {
			shutdownCh <- nil
		}
	}()

	// Listen for shutdown signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a shutdown signal, an error from the server, or a startup timeout.
	select {
	case <-quit:
		logger.Info("Shutdown signal received")
	case err := <-shutdownCh:
		if err != nil {
			logger.WithError(err).Error("Server error during startup or runtime")
			return fmt.Errorf("server error: %w", err)
		}
	case <-time.After(timeoutDuration):
		logger.Error("Server did not start within the expected timeout; startup timeout reached")
		return fmt.Errorf("server startup timeout after %s", timeoutDuration)
	}

	logger.Info("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Error during server shutdown")
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	elapsed := time.Since(startTime)
	logger.Infof("Server shutdown completed gracefully after %s", elapsed)
	return nil
}
