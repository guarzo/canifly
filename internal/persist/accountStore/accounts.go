// persist/accounts.go
package accountStore

import (
	"fmt"
	"github.com/guarzo/canifly/internal/persist"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const accountFileName = "accounts.json"

var _ interfaces.AccountRepository = (*AccountStore)(nil)

type AccountStore struct {
	logger interfaces.Logger
}

func NewAccountStore(logger interfaces.Logger) *AccountStore {
	as := &AccountStore{
		logger: logger,
	}

	return as
}

func (as *AccountStore) FetchAccounts() ([]model.Account, error) {
	filePath, err := as.getAccountFileName()
	if err != nil {
		return nil, err
	}
	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		as.logger.Info("No identity file found for accounts")
		return accounts, nil
	}

	if err := persist.OldReadJson(filePath, &accounts); err != nil {
		as.logger.WithError(err).Error("Error loading accounts")
		return nil, err
	}

	as.logger.Debugf("Loaded %d accounts", len(accounts))
	return accounts, nil
}

func (as *AccountStore) SaveAccounts(accounts []model.Account) error {
	filePath, err := as.getAccountFileName()
	if err != nil {
		return err
	}

	if err := persist.OldSaveJson(filePath, accounts); err != nil {
		as.logger.WithError(err).Error("Error saving accounts data")
		return fmt.Errorf("error saving accounts: %w", err)
	}

	as.logger.Debugf("Saved %d accounts", len(accounts))
	return nil
}

func (as *AccountStore) DeleteAccounts() error {
	idFile, err := as.getAccountFileName()
	if err != nil {
		return err
	}
	if err := os.Remove(idFile); err != nil && !os.IsNotExist(err) {
		as.logger.WithError(err).Errorf("Failed to delete account file %s", idFile)
		return err
	}
	as.logger.Info("account file deleted")
	return nil
}

func (as *AccountStore) getAccountFileName() (string, error) {
	accountPath, err := persist.GetWriteableSubPath(persist.AccountDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(accountPath, accountFileName), nil
}
