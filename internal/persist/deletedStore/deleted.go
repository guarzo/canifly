// persist/deleted.go
package deletedStore

import (
	"fmt"
	"github.com/guarzo/canifly/internal/persist"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

const deletedFileName = "deleted.json"

var _ interfaces.DeletedCharactersRepository = (*DeletedStore)(nil)

type DeletedStore struct {
	logger interfaces.Logger
}

func NewDeletedStore(l interfaces.Logger) *DeletedStore {
	return &DeletedStore{
		logger: l,
	}
}

func (ds *DeletedStore) SaveDeletedCharacters(chars []string) error {
	filename, err := getDeletedFileName()
	if err != nil {
		return fmt.Errorf("failed to save deleted characters")
	}

	if err := persist.OldSaveJson(filename, chars); err != nil {
		ds.logger.WithError(err).Errorf("Failed to save deleted characters to %s", filename)
		return err
	}

	return nil
}

func (ds *DeletedStore) FetchDeletedCharacters() ([]string, error) {
	filename, err := getDeletedFileName()
	if err != nil {
		return []string{}, fmt.Errorf("failed to load deleted characters")
	}

	var chars []string
	if err := persist.OldReadJson(filename, &chars); err != nil {
		if os.IsNotExist(err) {
			return []string{}, fmt.Errorf("deleted character file does not exist: %s", filename)
		}
		return []string{}, fmt.Errorf("failed to load deleted characters from %s, %v", filename, err)
	}

	return chars, nil
}

func getDeletedFileName() (string, error) {
	accountPath, err := persist.GetWriteableSubPath(persist.AccountDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(accountPath, deletedFileName), nil
}
