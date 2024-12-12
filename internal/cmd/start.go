package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/crypto"
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

	initializeComponents(cfg, logger)

	services, err := server.GetServices(logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	r := server.SetupHandlers(cfg.SecretKey, logger, services)
	return startServer(r, cfg.Port, logger)
}

func initializeComponents(cfg server.Config, logger interfaces.Logger) {
	key, err := base64.StdEncoding.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.WithError(err).Fatal("Failed to decode secret key.")
	}
	if err = crypto.Initialize(key); err != nil {
		logger.WithError(err).Fatal("Failed to initialize encryption.")
	}

	clientID, clientSecret, callbackURL := os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_CLIENT_SECRET"), os.Getenv("EVE_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		logger.Fatal("EVE_CLIENT_ID, EVE_CLIENT_SECRET, and EVE_CALLBACK_URL must be set")
	}

	auth.InitializeOAuth(clientID, clientSecret, callbackURL)

}

func startServer(r http.Handler, port string, logger interfaces.Logger) error {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

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

	logger.Infof("Server listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen: %w", err)
	}

	logger.Info("Server exited cleanly")
	return nil
}
