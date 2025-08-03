package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type testLogger struct{}

func (t testLogger) Debugf(format string, args ...interface{})                  {}
func (t testLogger) Debug(args ...interface{})                                  {}
func (t testLogger) Infof(format string, args ...interface{})                   {}
func (t testLogger) Info(args ...interface{})                                   {}
func (t testLogger) Warnf(format string, args ...interface{})                   {}
func (t testLogger) Warn(args ...interface{})                                   {}
func (t testLogger) Errorf(format string, args ...interface{})                  {}
func (t testLogger) Error(args ...interface{})                                  {}
func (t testLogger) Fatalf(format string, args ...interface{})                  {}
func (t testLogger) Fatal(args ...interface{})                                  {}
func (t testLogger) WithError(err error) interfaces.Logger                      { return t }
func (t testLogger) WithFields(fields map[string]interface{}) interfaces.Logger { return t }
func (t testLogger) WithField(key string, value interface{}) interfaces.Logger  { return t }

func TestLoadConfig_SecretKeyPersistence(t *testing.T) {
	// Save original env var
	origSecretKey := os.Getenv("SECRET_KEY")
	defer os.Setenv("SECRET_KEY", origSecretKey)

	// Clear SECRET_KEY env var to test persistence
	os.Unsetenv("SECRET_KEY")

	logger := testLogger{}

	// First load should create a new secret
	cfg1, err := LoadConfig(logger)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg1.SecretKey == "" {
		t.Error("Expected non-empty secret key")
	}

	// Verify secret file was created
	secretPath := filepath.Join(cfg1.BasePath, "secret.key")
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		t.Error("Secret key file was not created")
	}

	// Second load should return the same secret
	cfg2, err := LoadConfig(logger)
	if err != nil {
		t.Fatalf("Failed to load config second time: %v", err)
	}
	if cfg1.SecretKey != cfg2.SecretKey {
		t.Error("Expected same secret key on second load")
	}

	// Clean up
	os.RemoveAll(cfg1.BasePath)
}

func TestLoadConfig_EnvironmentOverride(t *testing.T) {
	// Save original env var
	origSecretKey := os.Getenv("SECRET_KEY")
	defer os.Setenv("SECRET_KEY", origSecretKey)

	// Set a specific secret key
	testSecret := "test-secret-key-12345"
	os.Setenv("SECRET_KEY", testSecret)

	logger := testLogger{}

	cfg, err := LoadConfig(logger)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.SecretKey != testSecret {
		t.Errorf("Expected secret key from env var, got %s", cfg.SecretKey)
	}

	// Clean up
	os.RemoveAll(cfg.BasePath)
}
