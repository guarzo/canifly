// services/interfaces/eve.go
package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
	"golang.org/x/oauth2"
)

type SkillService interface {
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillName(id int32) string
	GetSkillTypes() map[string]model.SkillType
	CheckIfDuplicatePlan(name string) bool
	ParseAndSaveSkillPlan(contents, name string) error
	GetSkillPlanFile(name string) ([]byte, error)
	DeleteSkillPlan(name string) error
	GetSkillTypeByID(id string) (model.SkillType, bool)
	GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string)
}

type SkillRepository interface {
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillPlanFile(name string) ([]byte, error)
	GetSkillTypes() map[string]model.SkillType
	SaveSkillPlan(planName string, skills map[string]model.Skill) error
	DeleteSkillPlan(planName string) error
	GetSkillTypeByID(id string) (model.SkillType, bool)
}

type EveProfilesService interface {
	LoadCharacterSettings() ([]model.EveProfile, error)
	BackupDir(targetDir, backupDir string) error

	SyncDir(subDir, charId, userId string) (int, int, error)
	SyncAllDir(baseSubDir, charId, userId string) (int, int, error)
}

type EveProfilesRepository interface {
	// ListSettingsFiles returns raw file info for character and user files in a given subdirectory of the settings directory.
	ListSettingsFiles(subDir, settingsDir string) ([]model.RawFileInfo, error)

	// BackupDirectory creates a tar.gz backup of all directories under targetDir that start with "settings_".
	BackupDirectory(targetDir, backupDir string) error

	// GetSubDirectories returns subdirectories in settingsDir that start with "settings_".
	GetSubDirectories(settingsDir string) ([]string, error)

	// SyncSubdirectory copies user and char file contents from one subdirectory to another to ensure all have consistent files.
	SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error)

	// SyncAllSubdirectories applies SyncSubdirectory logic to all subdirectories of settingsDir, using baseSubDir as the source.
	SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error)
}

type SystemRepository interface {
	GetSystemName(systemID int64) string
	LoadSystems() error
}

type ESIService interface {
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacter(id string) (*model.CharacterResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
	ResolveCharacterNames(charIds []string) (map[string]string, error)
	SaveEsiCache() error
	GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error)
	GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error)
}
