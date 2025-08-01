package interfaces

import "github.com/guarzo/canifly/internal/model"

// StorageService provides unified file storage operations
type StorageService interface {
	// Generic JSON operations
	LoadJSON(filename string, v interface{}) error
	SaveJSON(filename string, v interface{}) error

	// Account Data
	LoadAccountData() (*model.AccountData, error)
	SaveAccountData(data *model.AccountData) error
	DeleteAccountData() error

	// Config Data
	LoadConfigData() (*model.ConfigData, error)
	SaveConfigData(data *model.ConfigData) error

	// App State
	LoadAppState() (*model.AppState, error)
	SaveAppState(state *model.AppState) error

	// EVE Profiles
	LoadEveProfiles() (map[string]interface{}, error)
	SaveEveProfiles(data map[string]interface{}) error

	// Skill Plans
	LoadSkillPlan(filename string) ([]byte, error)
	SaveSkillPlan(filename string, content []byte) error
	ListSkillPlans() ([]string, error)

	// Cache
	LoadCache() (map[string][]byte, error)
	SaveCache(cache map[string][]byte) error
	
	// Deleted Characters
	LoadDeletedCharacters() ([]string, error)
	SaveDeletedCharacters(deleted []string) error
	
	// API Cache
	LoadAPICache() (map[string][]byte, error)
	SaveAPICache(cache map[string][]byte) error

	// Utilities
	EnsureDirectories() error
	Exists(filename string) bool
	Remove(filename string) error
	GetBasePath() string
}