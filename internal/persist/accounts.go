// persist/accounts.go
package persist

import (
	"fmt"
	"os"

	"github.com/guarzo/canifly/internal/model"
)

func (ds *DataStore) FetchAccounts() ([]model.Account, error) {
	filePath, err := ds.getAccountFileName()
	if err != nil {
		return nil, err
	}
	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		ds.logger.Info("No identity file found for accounts")
		return accounts, nil
	}

	if err := readJSONFromFile(filePath, &accounts); err != nil {
		ds.logger.WithError(err).Error("Error loading accounts")
		return nil, err
	}

	ds.logger.Debugf("Loaded %d accounts", len(accounts))
	return accounts, nil
}

func (ds *DataStore) SaveAccounts(accounts []model.Account) error {
	filePath, err := ds.getAccountFileName()
	if err != nil {
		return err
	}

	if err := saveJSONToFile(filePath, accounts); err != nil {
		ds.logger.WithError(err).Error("Error saving accounts data")
		return fmt.Errorf("error saving accounts: %w", err)
	}

	ds.logger.Debugf("Saved %d accounts", len(accounts))
	return nil
}

func (ds *DataStore) DeleteAccounts() error {
	idFile, err := ds.getAccountFileName()
	if err != nil {
		return err
	}
	if err := os.Remove(idFile); err != nil && !os.IsNotExist(err) {
		ds.logger.WithError(err).Errorf("Failed to delete account file %s", idFile)
		return err
	}
	ds.logger.Info("account file deleted")
	return nil
}

func (ds *DataStore) getAccountFileName() (string, error) {
	return getAccountFileName(accountFileName)
}
