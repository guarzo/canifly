package interfaces

import "github.com/guarzo/canifly/internal/model"

// ConfigurationService handles application configuration
type ConfigurationService interface {
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

}