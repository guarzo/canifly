// persist/identity.go
package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils"
)

func (ds *DataStore) FetchAccounts() ([]model.Account, error) {
	filePath := ds.getAccountFileName()
	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		ds.logger.Info("No identity file found for accounts")
		return accounts, nil
	}

	if err := ds.loadData(filePath, &accounts); err != nil {
		ds.logger.WithError(err).Error("Error loading accounts")
		return nil, err
	}

	ds.logger.Debugf("Loaded %d accounts", len(accounts))
	return accounts, nil
}

func (ds *DataStore) SaveAccounts(accounts []model.Account) error {
	filePath := ds.getAccountFileName()
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving accounts")
	}

	if err := ds.saveData(accounts, filePath); err != nil {
		ds.logger.WithError(err).Error("Error saving accounts data")
		return fmt.Errorf("error saving accounts: %w", err)
	}

	ds.logger.Debugf("Saved %d accounts", len(accounts))
	return nil
}

func (ds *DataStore) DeleteAccounts() error {
	idFile := ds.getAccountFileName()
	if err := os.Remove(idFile); err != nil && !os.IsNotExist(err) {
		ds.logger.WithError(err).Errorf("Failed to delete identity file %s", idFile)
		return err
	}
	ds.logger.Info("Identity file deleted")
	return nil
}

func (ds *DataStore) getAccountFileName() string {
	identityPath, err := ds.GetWriteablePath()
	if err != nil {
		ds.logger.WithError(err).Error("Error retrieving writable identity path")
		return ""
	}
	return filepath.Join(identityPath, "identity.json")
}

// saveData and loadData utilities
func (ds *DataStore) saveData(data interface{}, outputFile string) error {
	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to open file for writing: %s", outputFile)
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	if err := encoder.Encode(data); err != nil {
		ds.logger.WithError(err).Errorf("Failed to encode JSON to file: %s", outputFile)
		return err
	}

	return nil
}

func (ds *DataStore) loadData(inputFile string, data interface{}) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to open file: %s", inputFile)
		return err
	}
	defer inFile.Close()

	decoder := json.NewDecoder(inFile)
	if err := decoder.Decode(data); err != nil {
		ds.logger.WithError(err).Errorf("Failed to decode JSON from file: %s", inputFile)
		return err
	}

	return nil
}

// Encryption/Decryption wrappers (assuming utils package handles actual logic)
func (ds *DataStore) encryptData(data interface{}, outputFile string) error {
	ds.logger.Debugf("Encrypting data to %s", outputFile)
	return utils.EncryptData(data, outputFile)
}

func (ds *DataStore) decryptData(inputFile string, data interface{}) error {
	ds.logger.Debugf("Decrypting data from %s", inputFile)
	return utils.DecryptData(inputFile, data)
}
