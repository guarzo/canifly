package interfaces

import "github.com/guarzo/canifly/internal/model"

// ConfigurationService consolidates ConfigService, AppStateService and DashboardService interfaces
type ConfigurationService interface {
	// App State Management (from AppStateService)
	GetAppState() model.AppState
	SetAppStateLogin(isLoggedIn bool) error
	UpdateAndSaveAppState(data model.AppState) error
	ClearAppState()

	// Configuration Management (from ConfigService)
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
	SaveRoles(roles []string) error

	// Dashboard Management (from DashboardService)
	RefreshAccountsAndState() (model.AppState, error)
	RefreshDataInBackground() error
	GetCurrentAppState() model.AppState
}