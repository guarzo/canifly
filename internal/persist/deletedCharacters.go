package persist

import (
	"fmt"
	"os"
	"path/filepath"
)

func (ds *DataStore) getDeletedFileName() string {
	writePath, err := ds.GetWriteablePath()
	if err != nil {
		ds.logger.WithError(err).Error("Error retrieving writable data path for config")
		return ""
	}
	return filepath.Join(writePath, "deleted_characters.json")
}

func (ds *DataStore) SaveDeletedCharacters(chars []string) error {
	filename := ds.getDeletedFileName()
	if filename == "" {
		return fmt.Errorf("failed to save deleted characters")
	}

	if err := writeJSONToFile(filename, chars); err != nil {
		ds.logger.WithError(err).Errorf("Failed to save cache to %s", filename)
		return err
	}

	return nil
}

func (ds *DataStore) FetchDeletedCharacters() ([]string, error) {
	filename := ds.getDeletedFileName()
	if filename == "" {
		return []string{}, fmt.Errorf("failed to load deleted characters")
	}

	var chars []string
	if err := readJSONFromFile(filename, &chars); err != nil {
		if os.IsNotExist(err) {
			return []string{}, fmt.Errorf("deleted character file does not exist: %s", filename)
		}
		return []string{}, fmt.Errorf("failed to load deleted characters from %s, %v", filename, err)
	}

	return chars, nil
}
