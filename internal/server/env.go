package server

import (
	"encoding/base64"
	"fmt"
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
	// Try local .env first (for development)
	if err := godotenv.Load(); err != nil {
		logger.Debug("No local .env file found")
	}

	cfg := Config{}

	cfg.Port = getPort()
	
	secretKey, err := getSecretKey(logger)
	if err != nil {
		return cfg, err
	}
	cfg.SecretKey = secretKey

	// Get EVE credentials from environment first (takes precedence)
	cfg.ClientID = os.Getenv("EVE_CLIENT_ID")
	cfg.ClientSecret = os.Getenv("EVE_CLIENT_SECRET")
	cfg.CallbackURL = os.Getenv("EVE_CALLBACK_URL")

	// Note: We'll check stored credentials later after services are initialized
	// For now, we'll allow empty credentials and handle them in the handlers

	cfg.PathSuffix = os.Getenv("PATH_SUFFIX")
	configDir, err := os.UserConfigDir()
	if err != nil {
		return cfg, fmt.Errorf("unable to get user config dir: %v", err)
	}
	cfg.BasePath = filepath.Join(configDir, "canifly")
	
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		// Try fallback to temp directory
		cfg.BasePath = filepath.Join(os.TempDir(), "canifly")
		logger.Warnf("Failed to create config directory, using temp: %v", err)
		if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
			return cfg, fmt.Errorf("failed to create base directory: %w", err)
		}
	}

	return cfg, nil
}

// getSecretKey retrieves or generates the encryption secret key
func getSecretKey(logger interfaces.Logger) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		key, err := persist.GenerateSecret()
		if err != nil {
			return "", fmt.Errorf("failed to generate secret key: %w", err)
		}
		secret = base64.StdEncoding.EncodeToString(key)
		logger.Warn("Using a generated key for testing only.")
	}
	return secret, nil
}

// getPort returns the port the server should listen on
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "42423"
	}
	return port
}

// ValidateStartupPaths ensures all required directories exist and are writable
func ValidateStartupPaths(cfg Config, logger interfaces.Logger) error {
	// Check base path exists and is writable
	testFile := filepath.Join(cfg.BasePath, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("base path %s is not writable: %w", cfg.BasePath, err)
	}
	os.Remove(testFile)
	
	// Create subdirectories that services will need
	requiredDirs := []string{
		filepath.Join(cfg.BasePath, "plans"),
		filepath.Join(cfg.BasePath, "eve"),
	}
	
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Warnf("Failed to create directory %s: %v", dir, err)
		}
	}
	
	logger.Infof("Startup validation complete. Base path: %s", cfg.BasePath)
	return nil
}
