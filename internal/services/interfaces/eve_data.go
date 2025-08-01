package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
	"golang.org/x/oauth2"
)

// EVEDataService consolidates ESIService, CacheService, CharacterService, SkillService, and EveProfilesService
type EVEDataService interface {
	// ESI API Operations (from ESIService)
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacter(id string) (*model.CharacterResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
	ResolveCharacterNames(charIds []string) (map[string]string, error)
	GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error)
	GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error)

	// Character Management (from CharacterService)
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
	DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
	UpdateCharacterFields(characterID int64, updates map[string]interface{}) error
	RemoveCharacter(characterID int64) error

	// Skill Plan Management (from SkillService)
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillName(id int32) string
	GetSkillTypes() map[string]model.SkillType
	CheckIfDuplicatePlan(name string) bool
	ParseAndSaveSkillPlan(contents, name string) error
	GetSkillPlanFile(name string) ([]byte, error)
	DeleteSkillPlan(name string) error
	GetSkillTypeByID(id string) (model.SkillType, bool)
	GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string)

	// Eve Profile Management (from EveProfilesService)
	LoadCharacterSettings() ([]model.EveProfile, error)
	BackupDir(targetDir, backupDir string) error
	SyncDir(subDir, charId, userId string) (int, int, error)
	SyncAllDir(baseSubDir, charId, userId string) (int, int, error)

	// Cache Management
	SaveCache() error
	LoadCache() error
	SaveEsiCache() error
}