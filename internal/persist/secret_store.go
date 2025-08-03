package persist

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

const (
	secretKeyFile = "secret.key"
	keySize       = 32 // 256 bits
)

// SecretStore manages the application's encryption secret key
type SecretStore struct {
	basePath string
	fs       FileSystem
}

// NewSecretStore creates a new secret store
func NewSecretStore(basePath string, fs FileSystem) *SecretStore {
	return &SecretStore{
		basePath: basePath,
		fs:       fs,
	}
}

// GetOrCreateSecret retrieves the existing secret key or creates a new one
func (s *SecretStore) GetOrCreateSecret() (string, error) {
	keyPath := filepath.Join(s.basePath, secretKeyFile)

	// Try to read existing key
	data, err := s.fs.ReadFile(keyPath)
	if err == nil && len(data) > 0 {
		return string(data), nil
	}

	// If file doesn't exist or is empty, generate new key
	if os.IsNotExist(err) || len(data) == 0 {
		key, err := s.generateAndSaveKey(keyPath)
		if err != nil {
			return "", fmt.Errorf("failed to generate secret key: %w", err)
		}
		return key, nil
	}

	return "", fmt.Errorf("failed to read secret key: %w", err)
}

// generateAndSaveKey creates a new secret key and saves it
func (s *SecretStore) generateAndSaveKey(keyPath string) (string, error) {
	// Generate random bytes
	keyBytes := make([]byte, keySize)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Encode to base64 for storage
	encodedKey := base64.StdEncoding.EncodeToString(keyBytes)

	// Ensure directory exists
	dir := filepath.Dir(keyPath)
	if err := s.fs.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write key with restricted permissions (owner read/write only)
	if err := AtomicWriteFile(s.fs, keyPath, []byte(encodedKey), 0600); err != nil {
		return "", fmt.Errorf("failed to save secret key: %w", err)
	}

	return encodedKey, nil
}

// DeleteSecret removes the stored secret key (useful for testing or reset)
func (s *SecretStore) DeleteSecret() error {
	keyPath := filepath.Join(s.basePath, secretKeyFile)
	err := s.fs.Remove(keyPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete secret key: %w", err)
	}
	return nil
}
