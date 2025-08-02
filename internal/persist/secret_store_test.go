package persist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecretStore_GetOrCreateSecret(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "secret-store-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fs := NewOSFileSystem()
	store := NewSecretStore(tmpDir, fs)

	// First call should create a new secret
	secret1, err := store.GetOrCreateSecret()
	if err != nil {
		t.Fatalf("Failed to get or create secret: %v", err)
	}
	if secret1 == "" {
		t.Error("Expected non-empty secret")
	}

	// Verify the file was created
	keyPath := filepath.Join(tmpDir, secretKeyFile)
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("Secret key file was not created")
	}

	// Second call should return the same secret
	secret2, err := store.GetOrCreateSecret()
	if err != nil {
		t.Fatalf("Failed to get existing secret: %v", err)
	}
	if secret1 != secret2 {
		t.Error("Expected same secret on second call")
	}

	// Create a new store instance to simulate app restart
	store2 := NewSecretStore(tmpDir, fs)
	secret3, err := store2.GetOrCreateSecret()
	if err != nil {
		t.Fatalf("Failed to get secret after restart: %v", err)
	}
	if secret1 != secret3 {
		t.Error("Expected same secret after simulated restart")
	}
}

func TestSecretStore_DeleteSecret(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "secret-store-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fs := NewOSFileSystem()
	store := NewSecretStore(tmpDir, fs)

	// Create a secret
	_, err = store.GetOrCreateSecret()
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}

	// Delete the secret
	err = store.DeleteSecret()
	if err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	// Verify file is gone
	keyPath := filepath.Join(tmpDir, secretKeyFile)
	if _, err := os.Stat(keyPath); !os.IsNotExist(err) {
		t.Error("Secret key file was not deleted")
	}
}