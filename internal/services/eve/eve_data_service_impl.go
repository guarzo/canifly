package eve

import (
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	cacheSvc "github.com/guarzo/canifly/internal/services/cache"
	characterSvc "github.com/guarzo/canifly/internal/services/character"
	"github.com/guarzo/canifly/internal/services/interfaces"
	skillplanSvc "github.com/guarzo/canifly/internal/services/skillplan"
)

// Compile-time interface checks
var _ interfaces.ESIAPIService = (*EVEDataServiceImpl)(nil)
var _ interfaces.CharacterService = (*EVEDataServiceImpl)(nil)
var _ interfaces.SkillPlanService = (*EVEDataServiceImpl)(nil)
var _ interfaces.CacheableService = (*EVEDataServiceImpl)(nil)
var _ interfaces.EVEDataComposite = (*EVEDataServiceImpl)(nil)
var _ interfaces.EVEDataService = (*EVEDataServiceImpl)(nil)

// EVEDataServiceImpl is the full implementation without delegation
// It implements all the split interfaces: ESIAPIService, CharacterService,
// SkillPlanService, ProfileService, CacheableService, and EVEDataComposite
type EVEDataServiceImpl struct {
	// Core dependencies
	logger     interfaces.Logger
	httpClient interfaces.EsiHttpClient
	authClient interfaces.AuthClient
	configSvc  interfaces.ConfigurationService
	storage    interfaces.StorageService

	// Repositories still needed
	systemRepo interfaces.SystemRepository

	// Cache management is delegated to a standalone cache service
	persistentCache *cacheSvc.PersistentCacheService

	// Character methods are delegated to *character.Service
	character *characterSvc.Service

	// Skill plan methods are delegated to *skillplan.Service
	skillPlan *skillplanSvc.Service

	// ESI API methods are delegated to *ESIClient
	esiClient *ESIClient
}

// NewEVEDataServiceImpl creates a new consolidated EVE data service implementation
func NewEVEDataServiceImpl(
	logger interfaces.Logger,
	httpClient interfaces.EsiHttpClient,
	authClient interfaces.AuthClient,
	configSvc interfaces.ConfigurationService,
	storage interfaces.StorageService,
	systemRepo interfaces.SystemRepository,
	persistentCache *cacheSvc.PersistentCacheService,
	character *characterSvc.Service,
	skillPlan *skillplanSvc.Service,
	esiClient *ESIClient,
) interfaces.EVEDataService {
	svc := &EVEDataServiceImpl{
		logger:          logger,
		httpClient:      httpClient,
		authClient:      authClient,
		configSvc:       configSvc,
		storage:         storage,
		systemRepo:      systemRepo,
		persistentCache: persistentCache,
		character:       character,
		skillPlan:       skillPlan,
		esiClient:       esiClient,
	}

	return svc
}

// ESI API Operations — delegated to *ESIClient
func (s *EVEDataServiceImpl) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	return s.esiClient.GetUserInfo(token)
}

func (s *EVEDataServiceImpl) GetCharacter(id string) (*model.CharacterResponse, error) {
	return s.esiClient.GetCharacter(id)
}

func (s *EVEDataServiceImpl) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	return s.esiClient.GetCharacterSkills(characterID, token)
}

func (s *EVEDataServiceImpl) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	return s.esiClient.GetCharacterSkillQueue(characterID, token)
}

func (s *EVEDataServiceImpl) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	return s.esiClient.GetCharacterLocation(characterID, token)
}

func (s *EVEDataServiceImpl) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	return s.esiClient.ResolveCharacterNames(charIds)
}

func (s *EVEDataServiceImpl) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	return s.esiClient.GetCorporation(id, token)
}

func (s *EVEDataServiceImpl) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	return s.esiClient.GetAlliance(id, token)
}

// Character Management — delegated to *character.Service
func (s *EVEDataServiceImpl) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	return s.character.ProcessIdentity(charIdentity)
}

func (s *EVEDataServiceImpl) DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error) {
	return s.character.DoesCharacterExist(characterID)
}

func (s *EVEDataServiceImpl) UpdateCharacter(characterID int64, update model.CharacterUpdate) error {
	return s.character.UpdateCharacter(characterID, update)
}

func (s *EVEDataServiceImpl) RemoveCharacter(characterID int64) error {
	return s.character.RemoveCharacter(characterID)
}

func (s *EVEDataServiceImpl) RefreshCharacterData(characterID int64) (bool, error) {
	return s.character.RefreshCharacterData(characterID)
}

// Skill Plan Management — delegated to *skillplan.Service
func (s *EVEDataServiceImpl) GetSkillPlans() map[string]model.SkillPlan {
	return s.skillPlan.GetSkillPlans()
}

func (s *EVEDataServiceImpl) GetSkillName(id int32) string { return s.skillPlan.GetSkillName(id) }

func (s *EVEDataServiceImpl) GetSkillTypes() map[string]model.SkillType {
	return s.skillPlan.GetSkillTypes()
}

func (s *EVEDataServiceImpl) CheckIfDuplicatePlan(name string) bool {
	return s.skillPlan.CheckIfDuplicatePlan(name)
}

func (s *EVEDataServiceImpl) ParseAndSaveSkillPlan(contents, name string) error {
	return s.skillPlan.ParseAndSaveSkillPlan(contents, name)
}

func (s *EVEDataServiceImpl) GetSkillPlanFile(name string) ([]byte, error) {
	return s.skillPlan.GetSkillPlanFile(name)
}

func (s *EVEDataServiceImpl) DeleteSkillPlan(name string) error {
	return s.skillPlan.DeleteSkillPlan(name)
}

func (s *EVEDataServiceImpl) ListSkillPlans() ([]string, error) {
	return s.skillPlan.ListSkillPlans()
}

func (s *EVEDataServiceImpl) RefreshRemotePlans() error { return s.skillPlan.RefreshRemotePlans() }

func (s *EVEDataServiceImpl) GetSkillTypeByID(id string) (model.SkillType, bool) {
	return s.skillPlan.GetSkillTypeByID(id)
}

func (s *EVEDataServiceImpl) GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string) {
	return s.skillPlan.GetPlanAndConversionData(accounts, skillPlans, skillTypes)
}

// Cache Management — delegated to persistentCache
func (s *EVEDataServiceImpl) SaveCache() error    { return s.persistentCache.SaveCache() }
func (s *EVEDataServiceImpl) LoadCache() error    { return s.persistentCache.LoadCache() }
func (s *EVEDataServiceImpl) SaveEsiCache() error { return s.persistentCache.SaveEsiCache() }

func (s *EVEDataServiceImpl) Get(key string) ([]byte, bool) { return s.persistentCache.Get(key) }
func (s *EVEDataServiceImpl) Set(key string, value []byte, expiration time.Duration) {
	s.persistentCache.Set(key, value, expiration)
}

// Shutdown gracefully shuts down the EVE data service
func (s *EVEDataServiceImpl) Shutdown() { s.persistentCache.Shutdown() }

// Clear removes all cache entries
func (s *EVEDataServiceImpl) Clear() { s.persistentCache.Clear() }
