package interfaces

import (
	"time"

	"github.com/guarzo/canifly/internal/model"
)

// AccountDataRepository interface
type AccountDataRepository interface {
	// FetchAccountData retrieves the entire account domain data (Accounts, UserAccount map, and Associations).
	FetchAccountData() (model.AccountData, error)

	// SaveAccountData saves the entire account domain data structure.
	SaveAccountData(data model.AccountData) error

	// DeleteAccountData removes the persisted account data file, if any.
	DeleteAccountData() error

	// FetchAccounts returns only the Accounts slice from the account data.
	// This is a convenience method that internally fetches AccountData and returns AccountData.Accounts.
	FetchAccounts() ([]model.Account, error)

	// SaveAccounts updates the Accounts slice in the account data, leaving UserAccount and Associations unchanged.
	SaveAccounts(accounts []model.Account) error

	// DeleteAccounts clears out the Accounts slice (but not necessarily UserAccount or Associations).
	DeleteAccounts() error
}

// CacheRepository interface
type CacheRepository interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadApiCache() error
	SaveApiCache() error
}

// DeletedCharactersRepository interface
type DeletedCharactersRepository interface {
	FetchDeletedCharacters() ([]string, error)
	SaveDeletedCharacters([]string) error
}

// LoginRepository interface
type LoginRepository interface {
	Set(state string, authStatus *model.AuthStatus)
	Get(state string) (*model.AuthStatus, bool)
	Delete(state string)
}

// AppStateRepository interface
type AppStateRepository interface {
	// GetAppState returns the current in-memory AppState.
	GetAppState() model.AppState

	// SetAppState replaces the current AppState in memory.
	SetAppState(appState model.AppState)

	// SetAppStateLogin updates the LoggedIn field in AppState and persists the change.
	SetAppStateLogin(isLoggedIn bool) error

	// ClearAppState resets the AppState to an empty struct.
	ClearAppState()

	// SaveAppStateSnapshot writes the current AppState to disk.
	SaveAppStateSnapshot(appState model.AppState) error
}

// ConfigRepository interface
type ConfigRepository interface {
	// FetchConfigData loads the entire ConfigData structure.
	FetchConfigData() (*model.ConfigData, error)

	// SaveConfigData persists the entire ConfigData structure.
	SaveConfigData(*model.ConfigData) error

	// FetchUserSelections retrieves the UserSelections from the config data.
	FetchUserSelections() (model.DropDownSelections, error)

	// SaveUserSelections updates UserSelections in the config data.
	SaveUserSelections(selections model.DropDownSelections) error

	// FetchRoles retrieves the Roles slice from the config data.
	FetchRoles() ([]string, error)

	// SaveRoles updates the Roles slice in the config data.
	SaveRoles(roles []string) error

	GetDefaultSettingsDir() (string, error)
	BackupJSONFiles(backupDir string) error
}