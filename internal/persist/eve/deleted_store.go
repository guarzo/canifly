package eve

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	deletedFileName = "deleted.json"
)

var _ interfaces.DeletedCharactersRepository = (*DeletedStore)(nil)

type DeletedStore struct {
	logger      interfaces.Logger
	fs          persist.FileSystem
	basePath    string
	mu          sync.RWMutex
	cachedChars []string
}

func NewDeletedStore(l interfaces.Logger, fs persist.FileSystem, basePath string) *DeletedStore {
	return &DeletedStore{
		logger:   l,
		fs:       fs,
		basePath: basePath,
	}
}

func (ds *DeletedStore) SaveDeletedCharacters(chars []string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	filename := filepath.Join(ds.basePath, deletedFileName)
	if err := persist.SaveJsonToFile(ds.fs, filename, chars); err != nil {
		ds.logger.WithError(err).Errorf("Failed to save deleted characters to %s", filename)
		return err
	}

	// Update the in-memory cache
	ds.cachedChars = chars
	return nil
}

func (ds *DeletedStore) FetchDeletedCharacters() ([]string, error) {
	ds.mu.RLock()
	// If we have cached data, return a copy directly
	if ds.cachedChars != nil {
		defer ds.mu.RUnlock()
		return ds.copyChars(ds.cachedChars), nil
	}
	ds.mu.RUnlock()

	// If no cached data, load from disk
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Check again if cachedChars got populated by another goroutine in between locks
	if ds.cachedChars != nil {
		return ds.copyChars(ds.cachedChars), nil
	}

	filename := filepath.Join(ds.basePath, deletedFileName)
	if _, err := ds.fs.Stat(filename); os.IsNotExist(err) {
		// It's not necessarily an error if the file doesn't exist,
		// we can treat it as empty since we now have a cache.
		ds.logger.Infof("Deleted characters file not found: %s", filename)
		ds.cachedChars = []string{}
		return []string{}, nil
	} else if err != nil {
		return []string{}, fmt.Errorf("failed to stat deleted character file: %w", err)
	}

	var chars []string
	if err := persist.ReadJsonFromFile(ds.fs, filename, &chars); err != nil {
		return []string{}, fmt.Errorf("failed to load deleted characters from %s: %w", filename, err)
	}

	// Cache the loaded data
	ds.cachedChars = chars
	return ds.copyChars(chars), nil
}

// copyChars creates a copy of the slice to avoid sharing the underlying array.
func (ds *DeletedStore) copyChars(chars []string) []string {
	cpy := make([]string, len(chars))
	copy(cpy, chars)
	return cpy
}
