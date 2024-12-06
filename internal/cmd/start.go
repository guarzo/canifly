package cmd

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/server"
	"github.com/guarzo/canifly/internal/utils"
)

var version = "0.0.40"

func Start() {
	logger := server.SetupLogger()
	logger.Infof("Starting application, version %s", version)

	server.LoadEnv(logger)

	// Initialize application components
	secret := server.GetSecretKey(logger)
	initializeComponents(secret, logger)

	// Set up router and server
	r := server.SetupRouter(secret, logger)
	port := server.GetPort()
	startServer(r, port, logger)
}

// initializeComponents sets up the app's core components (crypto, OAuth, skill plans, skill types).
func initializeComponents(secret string, logger *logrus.Logger) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		logger.WithError(err).Fatal("Failed to decode secret key.")
	}
	if err = utils.Initialize(key); err != nil {
		logger.WithError(err).Fatal("Failed to initialize encryption.")
	}

	clientID, clientSecret, callbackURL := os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_CLIENT_SECRET"), os.Getenv("EVE_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		logger.Fatal("EVE_CLIENT_ID, EVE_CLIENT_SECRET, and EVE_CALLBACK_URL must be set")
	}
	api.InitializeOAuth(clientID, clientSecret, callbackURL)

	if err = persist.ProcessSkillPlans(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill plans.")
	}

	if err = persist.LoadSkillTypes(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill types.")
	}

	if err = embed.LoadStatic(); err != nil {
		logger.WithError(err).Fatal("Failed to load templates.")
	}

	writeable, err := persist.GetWritablePlansPath()
	logger.WithField("path", writeable).Info("Writable path initialized.")
}

// startServer starts the HTTP server with CORS middleware
func startServer(r *mux.Router, port string, logger *logrus.Logger) {
	logger.Infof("Listening on port %s", port)
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsOptions)(r)))
}
