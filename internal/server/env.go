package server

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type Config struct {
	Port              string
	SecretKey         string
	ClientID          string
	ClientSecret      string
	CallbackURL       string
	PathSuffix        string
	BasePath          string
	SkillPlansRepoURL string
}

// LoadConfig loads the application configuration.
// Configuration priority:
// 1. Environment variables (for development/override)
// 2. Stored configuration (loaded later by services)
// 3. Default values
func LoadConfig(logger interfaces.Logger) (Config, error) {
	cfg := Config{}

	cfg.Port = getPort()
	
	// Get base path first so we can use it for secret key storage
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
	
	// Get or create persistent secret key
	secretKey, err := getSecretKey(cfg.BasePath, logger)
	if err != nil {
		return cfg, err
	}
	cfg.SecretKey = secretKey

	// EVE OAuth2 credentials - environment variables take precedence over stored config
	// This allows developers to override stored credentials during development
	cfg.ClientID = os.Getenv("EVE_CLIENT_ID")
	cfg.ClientSecret = os.Getenv("EVE_CLIENT_SECRET")
	cfg.CallbackURL = os.Getenv("EVE_CALLBACK_URL")

	// Optional: Custom skill plans repository URL (defaults to guarzo/eve-skills)
	cfg.SkillPlansRepoURL = os.Getenv("SKILLPLANS_REPO_URL")

	// Note: If EVE credentials are not set via environment variables,
	// the config service will check for stored credentials and prompt
	// for first-run configuration if needed

	cfg.PathSuffix = os.Getenv("PATH_SUFFIX")

	return cfg, nil
}

// getSecretKey retrieves or generates the encryption secret key
func getSecretKey(basePath string, logger interfaces.Logger) (string, error) {
	// First check environment variable (for development/override)
	secret := os.Getenv("SECRET_KEY")
	if secret != "" {
		logger.Debug("Using SECRET_KEY from environment variable")
		return secret, nil
	}
	
	// Use persistent secret store
	fs := persist.NewOSFileSystem()
	secretStore := persist.NewSecretStore(basePath, fs)
	
	secret, err := secretStore.GetOrCreateSecret()
	if err != nil {
		return "", fmt.Errorf("failed to get or create secret key: %w", err)
	}
	
	logger.Info("Using persistent secret key from config directory")
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
		filepath.Join(cfg.BasePath, "plans"),        // Skill plans
		filepath.Join(cfg.BasePath, "eve"),          // EVE character data
		filepath.Join(cfg.BasePath, "config"),       // Application config
		filepath.Join(cfg.BasePath, "config", "fuzzworks"), // Fuzzworks static data
	}
	
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Warnf("Failed to create directory %s: %v", dir, err)
		}
	}
	
	logger.Infof("Startup validation complete. Base path: %s", cfg.BasePath)
	return nil
}
