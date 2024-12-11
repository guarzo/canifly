// persist/deleted.go
package persist

import (
	"fmt"
	"os"
)

func (ds *DataStore) getDeletedFileName() (string, error) {
	return getAccountFileName(deletedFileName)
}

func (ds *DataStore) SaveDeletedCharacters(chars []string) error {
	filename, err := ds.getDeletedFileName()
	if err != nil {
		return fmt.Errorf("failed to save deleted characters")
	}

	if err := saveJSONToFile(filename, chars); err != nil {
		ds.logger.WithError(err).Errorf("Failed to save deleted characters to %s", filename)
		return err
	}

	return nil
}

func (ds *DataStore) FetchDeletedCharacters() ([]string, error) {
	filename, err := ds.getDeletedFileName()
	if err != nil {
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
