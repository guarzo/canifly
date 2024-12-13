// account_store.go
package account

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	accountFileName = "account_data.json"
)

var _ interfaces.AccountDataRepository = (*AccountDataStore)(nil)

type AccountDataStore struct {
	logger     interfaces.Logger
	fs         persist.FileSystem
	basePath   string
	mu         sync.RWMutex
	cachedData *model.AccountData // In-memory cached data
}

func NewAccountDataStore(logger interfaces.Logger, fs persist.FileSystem, basePath string) *AccountDataStore {
	return &AccountDataStore{
		logger:   logger,
		fs:       fs,
		basePath: basePath,
	}
}

func (as *AccountDataStore) FetchAccounts() ([]model.Account, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	accountData, err := as.fetchAccountDataLocked()
	if err != nil {
		return nil, err
	}

	if accountData.Accounts == nil {
		as.logger.Infof("empty accounts retrieved")
		accountData.Accounts = make([]model.Account, 0)
	}
	return accountData.Accounts, nil
}

func (as *AccountDataStore) SaveAccounts(accounts []model.Account) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	accountData, err := as.fetchAccountDataLocked()
	if err != nil {
		return err
	}
	accountData.Accounts = accounts
	return as.saveAccountDataLocked(*accountData)
}

func (as *AccountDataStore) DeleteAccounts() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	accountData, err := as.fetchAccountDataLocked()
	if err != nil {
		return err
	}
	accountData.Accounts = []model.Account{}
	return as.saveAccountDataLocked(*accountData)
}

func (as *AccountDataStore) FetchAccountData() (model.AccountData, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	data, err := as.fetchAccountDataLocked()
	if err != nil {
		return model.AccountData{}, err
	}
	return *data, nil
}

func (as *AccountDataStore) SaveAccountData(data model.AccountData) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.saveAccountDataLocked(data)
}

func (as *AccountDataStore) DeleteAccountData() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	filePath := filepath.Join(as.basePath, accountFileName)
	if err := as.fs.Remove(filePath); err != nil && !os.IsNotExist(err) {
		as.logger.WithError(err).Errorf("Failed to delete account data file %s", filePath)
		return err
	}
	as.logger.Info("account data file deleted")

	// Clear the cache since we no longer have data on disk
	as.cachedData = nil
	return nil
}

// internal methods that assume locks are held

func (as *AccountDataStore) fetchAccountDataLocked() (*model.AccountData, error) {
	// If we have cached data, return it directly
	if as.cachedData != nil {
		return as.cachedData, nil
	}

	filePath := filepath.Join(as.basePath, accountFileName)
	fileInfo, err := as.fs.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		as.logger.Info("No account data file found")
		as.cachedData = &model.AccountData{
			Accounts:     []model.Account{},
			Associations: []model.Association{},
		}
		return as.cachedData, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat account data file: %w", err)
	}

	var data model.AccountData
	if err := persist.ReadJsonFromFile(as.fs, filePath, &data); err != nil {
		as.logger.WithError(err).Error("Error loading account data")
		return nil, err
	}

	as.logger.Debugf("Loaded account data with %d accounts", len(data.Accounts))
	// Cache the loaded data
	as.cachedData = &data
	return as.cachedData, nil
}

func (as *AccountDataStore) saveAccountDataLocked(data model.AccountData) error {
	filePath := filepath.Join(as.basePath, accountFileName)
	if err := persist.SaveJsonToFile(as.fs, filePath, data); err != nil {
		as.logger.WithError(err).Error("Error saving account data")
		return fmt.Errorf("error saving account data: %w", err)
	}

	// Update the cache with the newly saved data
	as.cachedData = &data
	as.logger.Debugf("Saved account data with %d accounts", len(data.Accounts))
	return nil
}
