// services/interfaces/config.go
package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
)

type DashboardService interface {
	RefreshAccountsAndState() (model.AppState, error)
	RefreshDataInBackground() error
	GetCurrentAppState() model.AppState
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithError(err error) Logger
	WithFields(fields map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
}
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

type AppStateService interface {
	GetAppState() model.AppState
	SetAppStateLogin(isLoggedIn bool) error
	UpdateAndSaveAppState(data model.AppState) error
	ClearAppState()
}

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

type ConfigService interface {
	UpdateSettingsDir(dir string) error
	GetSettingsDir() (string, error)
	EnsureSettingsDir() error
	SaveUserSelections(model.DropDownSelections) error
	FetchUserSelections() (model.DropDownSelections, error)
	UpdateRoles(newRole string) error
	GetRoles() ([]string, error)
	UpdateBackupDir(dir string) error
	BackupJSONFiles(backupDir string) error
	FetchConfigData() (*model.ConfigData, error)
}
