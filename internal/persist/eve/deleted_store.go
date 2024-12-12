package eve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	deletedFileName = "deleted.json"
)

var _ interfaces.DeletedCharactersRepository = (*DeletedStore)(nil)

type DeletedStore struct {
	logger   interfaces.Logger
	fs       persist.FileSystem
	basePath string
}

func NewDeletedStore(l interfaces.Logger, fs persist.FileSystem, basePath string) *DeletedStore {
	return &DeletedStore{
		logger:   l,
		fs:       fs,
		basePath: basePath,
	}
}

func (ds *DeletedStore) SaveDeletedCharacters(chars []string) error {
	filename := filepath.Join(ds.basePath, deletedFileName)

	if err := persist.SaveJsonToFile(ds.fs, filename, chars); err != nil {
		ds.logger.WithError(err).Errorf("Failed to save deleted characters to %s", filename)
		return err
	}

	return nil
}

func (ds *DeletedStore) FetchDeletedCharacters() ([]string, error) {
	filename := filepath.Join(ds.basePath, deletedFileName)

	if _, err := ds.fs.Stat(filename); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("deleted character file does not exist: %s", filename)
	} else if err != nil {
		return []string{}, fmt.Errorf("failed to stat deleted character file: %w", err)
	}

	var chars []string
	if err := persist.ReadJsonFromFile(ds.fs, filename, &chars); err != nil {
		return []string{}, fmt.Errorf("failed to load deleted characters from %s: %w", filename, err)
	}

	return chars, nil
}
