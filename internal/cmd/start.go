package cmd

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/guarzo/canifly/internal/auth"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/server"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/utils"
)

var version = "0.0.40"

func Start() {
	logger := server.SetupLogger()
	logger.Infof("Starting application, version %s", version)

	server.LoadEnv(logger)

	secret := server.GetSecretKey(logger)
	initializeComponents(secret, logger)

	httpClient := flyHttp.NewAPIClient("https://esi.evetech.net", "", logger)
	authClient := auth.NewAuthClient(logger)

	r := server.SetupRouter(secret, logger, server.GetServices(logger, authClient, httpClient))
	port := server.GetPort()
	startServer(r, port, logger)
}

func initializeComponents(secret string, logger interfaces.Logger) {
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

	auth.InitializeOAuth(clientID, clientSecret, callbackURL)

}

func startServer(r *mux.Router, port string, logger interfaces.Logger) {
	logger.Infof("Listening on port %s", port)
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(),
	)

	log.Fatal(http.ListenAndServe(":"+port, corsOptions(r)))
}
