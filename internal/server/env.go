package server

import (
	"encoding/base64"
	"fmt"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Config struct {
	Port         string
	SecretKey    string
	ClientID     string
	ClientSecret string
	CallbackURL  string
	PathSuffix   string
	BasePath     string
}

func LoadConfig(logger interfaces.Logger) (Config, error) {
	// Try local .env
	if err := godotenv.Load(); err != nil {
		// Try embedded .env
		embeddedEnv, err := embed.EnvFiles.Open("config/.env")
		if err != nil {
			logger.Warn("Failed to load embedded .env file. Using system environment variables.")
		} else {
			defer embeddedEnv.Close()
			envMap, err := godotenv.Parse(embeddedEnv)
			if err != nil {
				logger.WithError(err).Warn("Failed to parse embedded .env file.")
			} else {
				for key, value := range envMap {
					os.Setenv(key, value)
				}
			}
		}
	}

	cfg := Config{}

	cfg.Port = getPort()
	cfg.SecretKey = getSecretKey(logger)

	cfg.ClientID = os.Getenv("EVE_CLIENT_ID")
	cfg.ClientSecret = os.Getenv("EVE_CLIENT_SECRET")
	cfg.CallbackURL = os.Getenv("EVE_CALLBACK_URL")

	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.CallbackURL == "" {
		return cfg, fmt.Errorf("EVE_CLIENT_ID, EVE_CLIENT_SECRET, and EVE_CALLBACK_URL must be set")
	}

	cfg.PathSuffix = os.Getenv("PATH_SUFFIX")
	configDir, err := os.UserConfigDir()
	if err != nil {
		return cfg, fmt.Errorf("unable to get user config dir: %v", err)
	}
	cfg.BasePath = filepath.Join(configDir, "canifly")

	return cfg, nil
}

// getSecretKey retrieves or generates the encryption secret key
func getSecretKey(logger interfaces.Logger) string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		key, err := persist.GenerateSecret()
		if err != nil {
			logger.WithError(err).Fatal("Failed to generate secret key")
		}
		secret = base64.StdEncoding.EncodeToString(key)
		logger.Warn("Using a generated key for testing only.")
	}
	return secret
}

// getPort returns the port the server should listen on
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8713"
	}
	return port
}
