// services/interfaces/repositories.go
package interfaces

import (
	"time"

	"github.com/guarzo/canifly/internal/model"
)

type SettingsRepository interface {
	ConfigRepository
	UserSelectionsRepository
	FileSystemRepository
}

type ConfigRepository interface {
	FetchConfigData() (*model.ConfigData, error)
	SaveConfigData(*model.ConfigData) error
}

type LoginRepository interface {
	Set(state, value string)
	Get(state string) (string, bool)
	Delete(state string)
}

type AccountRepository interface {
	FetchAccounts() ([]model.Account, error)
	SaveAccounts([]model.Account) error
	DeleteAccounts() error
}

type UserSelectionsRepository interface {
	SaveUserSelections(model.UserSelections) error
	FetchUserSelections() (model.UserSelections, error)
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
	BackupDirectory(targetDir, backupDir string) error
	GetHomeDir() (string, error)
}

type StateRepository interface {
	SetAppState(as model.AppState)
	GetAppState() model.AppState
	SetAppStateLogin(in bool) error
	ClearAppState()
	SaveAppStateSnapshot(as model.AppState) error
}

type AssocRepository interface {
	ConfigRepository
	AccountRepository
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
	SaveSkillPlan(planName string, skills map[string]model.Skill) error
	DeleteSkillPlan(planName string) error
	GetWriteablePlansPath() (string, error)
	GetSkillTypeByID(id string) (model.SkillType, bool)
}
