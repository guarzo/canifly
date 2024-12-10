// services/interfaces/repositories.go
package interfaces

import (
	"time"

	"github.com/guarzo/canifly/internal/model"
)

type ConfigRepository interface {
	FetchConfigData() (*model.ConfigData, error)
	SaveConfigData(*model.ConfigData) error
}

type AccountRepository interface {
	FetchAccounts() ([]model.Account, error)
	SaveAccounts([]model.Account) error
	DeleteAccounts() error
}

type UserSelectionsRepository interface {
	SaveUserSelections(model.UserSelections) error
}

type DeletedCharactersRepository interface {
	FetchDeletedCharacters() ([]string, error)
	SaveDeletedCharacters([]string) error
}

type FileSystemRepository interface {
	GetSubDirectories(dir string) ([]string, error)
	GetFilesForDropdown(sd, settingsDir string) ([]model.CharFile, []model.UserFile, error)
	SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error)
	SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error)
	BackupDirectory(settingsDir string) (bool, string)
	GetHomeDir() (string, error)
}

type StateRepository interface {
	SetAppState(as model.AppState)
	GetAppState() model.AppState
	SetAppStateLogin(in bool) error
	ClearAppState() error
	SaveAppStateSnapshot(as model.AppState) error
}

type CacheRepository interface {
	GetFromCache(key string) ([]byte, bool)
	SetToCache(key string, value []byte, expiration time.Duration)
	LoadApiCache() error
	SaveApiCache() error
}

type SkillRepository interface {
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillTypes() map[string]model.SkillType
}
