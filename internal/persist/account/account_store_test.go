package account_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/account"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/stretchr/testify/assert"
)

// MockLogger that does nothing.
type MockLogger struct{}

func (m *MockLogger) Debug(args ...interface{})                                  {}
func (m *MockLogger) Debugf(format string, args ...interface{})                  {}
func (m *MockLogger) Info(args ...interface{})                                   {}
func (m *MockLogger) Infof(format string, args ...interface{})                   {}
func (m *MockLogger) Warn(args ...interface{})                                   {}
func (m *MockLogger) Warnf(format string, args ...interface{})                   {}
func (m *MockLogger) Error(args ...interface{})                                  {}
func (m *MockLogger) Errorf(format string, args ...interface{})                  {}
func (m *MockLogger) Fatal(args ...interface{})                                  {}
func (m *MockLogger) Fatalf(format string, args ...interface{})                  {}
func (m *MockLogger) WithError(err error) interfaces.Logger                      { return m }
func (m *MockLogger) WithField(key string, value interface{}) interfaces.Logger  { return m }
func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.Logger { return m }

func TestAccountDataStore_EmptyOnStart(t *testing.T) {
	logger := &MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := account.NewAccountDataStore(logger, fs, basePath)

	// Initially, no file exists, FetchAccountData should return empty data
	ad, err := store.FetchAccountData()
	assert.NoError(t, err)
	assert.Empty(t, ad.Accounts)
	assert.Empty(t, ad.Associations)

	// Similarly, FetchAccounts should return empty slice
	accounts, err := store.FetchAccounts()
	assert.NoError(t, err)
	assert.Empty(t, accounts)
}

func TestAccountDataStore_SaveAndFetch(t *testing.T) {
	logger := &MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := account.NewAccountDataStore(logger, fs, basePath)

	// Create some account data
	ad := model.AccountData{
		Accounts: []model.Account{
			{
				Name:       "TestAccount",
				Status:     "Alpha",
				Characters: []model.CharacterIdentity{},
				ID:         12345,
			},
		},
		Associations: []model.Association{
			{
				UserId:   "100",
				CharId:   "200",
				CharName: "TestChar",
			},
		},
	}

	// SaveAccountData
	err := store.SaveAccountData(ad)
	assert.NoError(t, err)

	// FetchAccountData and verify
	fetchedAd, err := store.FetchAccountData()
	assert.NoError(t, err)
	assert.Equal(t, ad, fetchedAd)

	// Test FetchAccounts convenience method
	accounts, err := store.FetchAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, "TestAccount", accounts[0].Name)
}

func TestAccountDataStore_SaveAccounts(t *testing.T) {
	logger := &MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := account.NewAccountDataStore(logger, fs, basePath)

	// Save some initial data
	err := store.SaveAccountData(model.AccountData{
		Accounts: []model.Account{
			{Name: "Acc1"},
		},
	})
	assert.NoError(t, err)

	// Now use SaveAccounts to overwrite just the Accounts
	newAccounts := []model.Account{
		{Name: "Acc2", ID: 999},
	}
	err = store.SaveAccounts(newAccounts)
	assert.NoError(t, err)

	// Verify that Accounts changed and Associations stayed empty
	ad, err := store.FetchAccountData()
	assert.NoError(t, err)
	assert.Len(t, ad.Accounts, 1)
	assert.Equal(t, "Acc2", ad.Accounts[0].Name)
	assert.Empty(t, ad.Associations)
}

func TestAccountDataStore_DeleteAccountData(t *testing.T) {
	logger := &MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := account.NewAccountDataStore(logger, fs, basePath)

	// Save some data
	err := store.SaveAccountData(model.AccountData{
		Accounts: []model.Account{
			{Name: "ToDelete"},
		},
	})
	assert.NoError(t, err)

	// Ensure file created
	filePath := filepath.Join(basePath, "account_data.json")
	_, err = fs.Stat(filePath)
	assert.NoError(t, err)

	// DeleteAccountData
	err = store.DeleteAccountData()
	assert.NoError(t, err)

	// Ensure file does not exist now
	_, err = fs.Stat(filePath)
	assert.True(t, os.IsNotExist(err))

	// FetchAccountData should return empty
	ad, err := store.FetchAccountData()
	assert.NoError(t, err)
	assert.Empty(t, ad.Accounts)
	assert.Empty(t, ad.Associations)
}

func TestAccountDataStore_DeleteAccounts(t *testing.T) {
	logger := &MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := account.NewAccountDataStore(logger, fs, basePath)

	// Save some data with Accounts
	err := store.SaveAccountData(model.AccountData{
		Accounts: []model.Account{
			{Name: "AccToDelete"},
		},
		Associations: []model.Association{
			{UserId: "U1", CharId: "C1"},
		},
	})
	assert.NoError(t, err)

	// Delete just the accounts
	err = store.DeleteAccounts()
	assert.NoError(t, err)

	// Verify accounts empty, but associations remain
	ad, err := store.FetchAccountData()
	assert.NoError(t, err)
	assert.Empty(t, ad.Accounts)
	assert.Len(t, ad.Associations, 1)
}
