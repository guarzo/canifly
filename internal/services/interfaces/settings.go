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

type SettingsRepository interface {
	ConfigRepository
	UserSelectionsRepository
	FileSystemRepository
}

type ConfigRepository interface {
	FetchConfigData() (*model.ConfigData, error)
	SaveConfigData(*model.ConfigData) error
}

type UserSelectionsRepository interface {
	SaveUserSelections(model.UserSelections) error
	FetchUserSelections() (model.UserSelections, error)
}

type FileSystemRepository interface {
	GetSubDirectories(dir string) ([]string, error)
	GetFilesForDropdown(sd, settingsDir string) ([]model.CharFile, []model.UserFile, error)
	SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error)
	SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error)
	BackupDirectory(targetDir, backupDir string) error
	GetHomeDir() (string, error)
}
