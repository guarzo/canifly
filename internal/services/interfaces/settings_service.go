package interfaces

import "github.com/guarzo/canifly/internal/model"

type SettingsService interface {
	UpdateSettingsDir(dir string) error
	GetSettingsDir() (string, error)
	EnsureSettingsDir() error
	LoadCharacterSettings() ([]model.SubDirData, error)
	BackupDir(targetDir, backupDir string) error
	SyncDir(subDir, charId, userId string) (int, int, error)
	SyncAllDir(baseSubDir, charId, userId string) (int, int, error)
	SaveUserSelections(model.UserSelections) error
	FetchUserSelections() (model.UserSelections, error)
	FetchConfigData() (*model.ConfigData, error)
	UpdateRoles(newRole string) error
}
