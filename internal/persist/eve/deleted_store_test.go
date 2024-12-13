package eve_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDeletedStore_FetchNoFile(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewDeletedStore(logger, fs, basePath)

	// With the updated store, no file means an empty slice and no error.
	chars, err := store.FetchDeletedCharacters()
	assert.NoError(t, err, "Should not error if file does not exist")
	assert.Empty(t, chars, "Should return empty slice if no file found")
}

func TestDeletedStore_SaveAndFetch(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewDeletedStore(logger, fs, basePath)

	original := []string{"char1", "char2", "char3"}

	// Save
	err := store.SaveDeletedCharacters(original)
	assert.NoError(t, err, "Save should succeed")

	// Fetch and verify
	fetched, err := store.FetchDeletedCharacters()
	assert.NoError(t, err, "Fetch should succeed after save")
	assert.Equal(t, original, fetched)
}

func TestDeletedStore_InvalidJSON(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewDeletedStore(logger, fs, basePath)

	// Write invalid JSON to file
	filePath := filepath.Join(basePath, "deleted.json")
	err := os.WriteFile(filePath, []byte("not valid json"), 0600)
	assert.NoError(t, err)

	// Fetch should fail
	chars, err := store.FetchDeletedCharacters()
	assert.Error(t, err)
	assert.Empty(t, chars)
	assert.Contains(t, err.Error(), "failed to load deleted characters")
}

func TestDeletedStore_StatError(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping permission test in CI environment")
	}

	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewDeletedStore(logger, fs, basePath)

	// Create a file but then revoke permissions on directory to cause stat error
	filePath := filepath.Join(basePath, "deleted.json")
	err := os.WriteFile(filePath, []byte(`["charX"]`), 0600)
	assert.NoError(t, err)

	// Revoke permissions so stat fails
	err = os.Chmod(basePath, 0000)
	assert.NoError(t, err, "Should be able to remove permissions")

	defer os.Chmod(basePath, 0755) // restore permissions for cleanup

	chars, err := store.FetchDeletedCharacters()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat deleted character file")
	assert.Empty(t, chars)
}
