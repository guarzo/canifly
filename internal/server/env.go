package server

import (
	"encoding/base64"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/utils"
	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from various sources
func LoadEnv(logger interfaces.Logger) {
	if err := godotenv.Load(); err != nil {
		embeddedEnv, err := embed.StaticFiles.Open("static/.env")
		if err != nil {
			logger.Warn("Failed to load embedded .env file. Using system environment variables.")
			return
		}
		defer embeddedEnv.Close()

		envMap, err := godotenv.Parse(embeddedEnv)
		if err != nil {
			logger.WithError(err).Warn("Failed to parse embedded .env file.")
			return
		}
		for key, value := range envMap {
			os.Setenv(key, value)
		}
	}
}

// GetSecretKey retrieves or generates the encryption secret key
func GetSecretKey(logger interfaces.Logger) string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		key, err := utils.GenerateSecret()
		if err != nil {
			logger.WithError(err).Fatal("Failed to generate secret key")
		}
		secret = base64.StdEncoding.EncodeToString(key)
		logger.Warn("Using a generated key for testing only.")
	}
	return secret
}

// GetPort returns the port the server should listen on
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8713"
	}
	return port
}

// GetWritablePath returns a writable path
func GetWritablePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, "canifly", "identity")
	pathSuffix := os.Getenv("PATH_SUFFIX")
	if pathSuffix != "" {
		path = filepath.Join(configDir, "canifly", pathSuffix)
	}
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	return path, nil
}
