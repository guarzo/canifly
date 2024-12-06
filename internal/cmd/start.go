package cmd

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/embed"
	httpServices "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/server"
	"github.com/guarzo/canifly/internal/services"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/utils"
)

var version = "0.0.40"

func Start() {
	logger := server.SetupLogger()
	logger.Infof("Starting application, version %s", version)

	server.LoadEnv(logger)

	secret := server.GetSecretKey(logger)
	initializeComponents(secret, logger)

	// Initialize dependencies
	httpClient := httpServices.NewAPIClient("https://esi.evetech.net", "", logger)
	authClient := auth.NewAuthClient()
	esiService := esi.NewESIService(httpClient, authClient)
	configService := services.NewConfigService(logger)
	skillService := services.NewSkillService(logger)

	r := server.SetupRouter(secret, logger, esiService, skillService, configService)
	port := server.GetPort()
	startServer(r, port, logger)
}

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

	auth.InitializeOAuth(clientID, clientSecret, callbackURL)

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

func startServer(r *mux.Router, port string, logger *logrus.Logger) {
	logger.Infof("Listening on port %s", port)
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsOptions)(r)))
}
