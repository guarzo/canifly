package eve

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	cacheSvc "github.com/guarzo/canifly/internal/services/cache"
	characterSvc "github.com/guarzo/canifly/internal/services/character"
	"github.com/guarzo/canifly/internal/services/interfaces"
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
	logger      interfaces.Logger
	httpClient  interfaces.EsiHttpClient
	authClient  interfaces.AuthClient
	accountMgmt interfaces.AccountManagementService
	configSvc   interfaces.ConfigurationService
	storage     interfaces.StorageService

	// Repositories still needed
	skillRepo  interfaces.SkillRepository
	systemRepo interfaces.SystemRepository

	// Cache management is delegated to a standalone cache service
	persistentCache *cacheSvc.PersistentCacheService

	// Character methods are delegated to *character.Service
	character *characterSvc.Service
}

// NewEVEDataServiceImpl creates a new consolidated EVE data service implementation
func NewEVEDataServiceImpl(
	logger interfaces.Logger,
	httpClient interfaces.EsiHttpClient,
	authClient interfaces.AuthClient,
	accountMgmt interfaces.AccountManagementService,
	configSvc interfaces.ConfigurationService,
	storage interfaces.StorageService,
	skillRepo interfaces.SkillRepository,
	systemRepo interfaces.SystemRepository,
	persistentCache *cacheSvc.PersistentCacheService,
	character *characterSvc.Service,
) interfaces.EVEDataService {
	svc := &EVEDataServiceImpl{
		logger:          logger,
		httpClient:      httpClient,
		authClient:      authClient,
		accountMgmt:     accountMgmt,
		configSvc:       configSvc,
		storage:         storage,
		skillRepo:       skillRepo,
		systemRepo:      systemRepo,
		persistentCache: persistentCache,
		character:       character,
	}

	return svc
}

// ESI API Operations — delegated to *character.Service
func (s *EVEDataServiceImpl) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	return s.character.GetUserInfo(token)
}

func (s *EVEDataServiceImpl) GetCharacter(id string) (*model.CharacterResponse, error) {
	return s.character.GetCharacter(id)
}

func (s *EVEDataServiceImpl) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	return s.character.GetCharacterSkills(characterID, token)
}

func (s *EVEDataServiceImpl) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	return s.character.GetCharacterSkillQueue(characterID, token)
}

func (s *EVEDataServiceImpl) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	return s.character.GetCharacterLocation(characterID, token)
}

func (s *EVEDataServiceImpl) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	return s.character.ResolveCharacterNames(charIds)
}

func (s *EVEDataServiceImpl) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	return s.character.GetCorporation(id, token)
}

func (s *EVEDataServiceImpl) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	return s.character.GetAlliance(id, token)
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

// Skill Plan Management
func (s *EVEDataServiceImpl) GetSkillPlans() map[string]model.SkillPlan {
	return s.skillRepo.GetSkillPlans()
}

func (s *EVEDataServiceImpl) GetSkillName(id int32) string {
	skillType, found := s.skillRepo.GetSkillTypeByID(strconv.Itoa(int(id)))
	if found {
		return skillType.TypeName
	}
	return ""
}

func (s *EVEDataServiceImpl) GetSkillTypes() map[string]model.SkillType {
	return s.skillRepo.GetSkillTypes()
}

func (s *EVEDataServiceImpl) CheckIfDuplicatePlan(name string) bool {
	plans := s.skillRepo.GetSkillPlans()
	_, exists := plans[name]
	return exists
}

func (s *EVEDataServiceImpl) ParseAndSaveSkillPlan(contents, name string) error {
	skills := s.parseSkillPlan(contents)
	return s.skillRepo.SaveSkillPlan(name, skills)
}

func (s *EVEDataServiceImpl) parseSkillPlan(contents string) map[string]model.Skill {
	skills := make(map[string]model.Skill)
	lines := strings.Split(contents, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Find the last space to separate skill name from level
		lastSpaceIndex := strings.LastIndex(line, " ")
		if lastSpaceIndex == -1 {
			s.logger.Warnf("No space found in line '%s'; skipping.", line)
			continue
		}

		skillName := line[:lastSpaceIndex]
		skillLevelStr := line[lastSpaceIndex+1:]

		// Parse skill level, handling Roman numerals if necessary
		skillLevel, err := parseSkillLevel(skillLevelStr)
		if err != nil {
			s.logger.Warnf("Invalid skill level '%s'; skipping line.", skillLevelStr)
			continue
		}

		// Check if the skill already exists and add/update if necessary
		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}

	return skills
}

var romanToInt = map[string]int{
	"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5,
}

func parseSkillLevel(levelStr string) (int, error) {
	if val, ok := romanToInt[levelStr]; ok {
		return val, nil
	}
	return strconv.Atoi(levelStr)
}

func (s *EVEDataServiceImpl) GetSkillPlanFile(name string) ([]byte, error) {
	return s.skillRepo.GetSkillPlanFile(name)
}

func (s *EVEDataServiceImpl) DeleteSkillPlan(name string) error {
	return s.skillRepo.DeleteSkillPlan(name)
}

func (s *EVEDataServiceImpl) ListSkillPlans() ([]string, error) {
	plans := s.GetSkillPlans()
	planNames := make([]string, 0, len(plans))
	for name := range plans {
		planNames = append(planNames, name)
	}
	return planNames, nil
}

func (s *EVEDataServiceImpl) RefreshRemotePlans() error {
	// Reload skill plans from repository (which will trigger GitHub download if configured)
	return s.skillRepo.LoadSkillPlans()
}

func (s *EVEDataServiceImpl) GetSkillTypeByID(id string) (model.SkillType, bool) {
	return s.skillRepo.GetSkillTypeByID(id)
}

func (s *EVEDataServiceImpl) GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string) {
	// Step 1: Initialize updatedSkillPlans and eveConversions
	updatedSkillPlans := s.initializeUpdatedPlans(skillPlans)
	eveConversions := s.initializeEveConversions(skillPlans, skillTypes)

	// Step 2: Process all accounts and characters
	typeIds := s.processAccountsAndCharacters(accounts, skillPlans, skillTypes, updatedSkillPlans)

	// Step 3: Convert skill IDs into names and update eveConversions
	s.updateEveConversionsWithSkillNames(typeIds, eveConversions)

	return updatedSkillPlans, eveConversions
}

func (s *EVEDataServiceImpl) initializeUpdatedPlans(skillPlans map[string]model.SkillPlan) map[string]model.SkillPlanWithStatus {
	updated := make(map[string]model.SkillPlanWithStatus)
	for planName, plan := range skillPlans {
		updated[planName] = model.SkillPlanWithStatus{
			Name:                plan.Name,
			Skills:              plan.Skills,
			QualifiedCharacters: []string{},
			PendingCharacters:   []string{},
			MissingCharacters:   []string{},
			MissingSkills:       make(map[string]map[string]int32),
			Characters:          []model.CharacterSkillPlanStatus{},
		}
	}
	return updated
}

func (s *EVEDataServiceImpl) initializeEveConversions(
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType,
) map[string]string {
	conversions := make(map[string]string)
	for planName := range skillPlans {
		if planType, exists := skillTypes[planName]; exists {
			conversions[planName] = planType.TypeID
		}
	}
	return conversions
}

func (s *EVEDataServiceImpl) processAccountsAndCharacters(
	accounts []model.Account,
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType,
	updatedSkillPlans map[string]model.SkillPlanWithStatus,
) []int32 {
	var typeIds []int32
	for _, account := range accounts {
		for _, chData := range account.Characters {
			character := chData.Character

			// Extract character skill and queue info
			characterSkills := s.mapCharacterSkills(character, &typeIds)
			skillQueueLevels := s.mapSkillQueueLevels(character)

			s.ensureCharacterMaps(&character)

			// Evaluate each plan for this character
			for planName, plan := range skillPlans {
				planResult := s.evaluatePlanForCharacter(plan, skillTypes, characterSkills, skillQueueLevels)

				planStatus := updatedSkillPlans[planName]
				s.updatePlanAndCharacterStatus(
					&planStatus,
					&character,
					planName,
					planResult,
				)
				updatedSkillPlans[planName] = planStatus
			}
		}
	}
	return typeIds
}

func (s *EVEDataServiceImpl) mapCharacterSkills(character model.Character, typeIds *[]int32) map[int32]int32 {
	skillsMap := make(map[int32]int32)
	for _, skill := range character.Skills {
		skillsMap[skill.SkillID] = skill.TrainedSkillLevel
		*typeIds = append(*typeIds, skill.SkillID)
	}
	return skillsMap
}

func (s *EVEDataServiceImpl) mapSkillQueueLevels(character model.Character) map[int32]struct {
	level      int32
	finishDate *time.Time
} {
	queueMap := make(map[int32]struct {
		level      int32
		finishDate *time.Time
	})
	for _, queuedSkill := range character.SkillQueue {
		current, exists := queueMap[queuedSkill.SkillID]
		if !exists || queuedSkill.FinishedLevel > current.level {
			queueMap[queuedSkill.SkillID] = struct {
				level      int32
				finishDate *time.Time
			}{level: queuedSkill.FinishedLevel, finishDate: queuedSkill.FinishDate}
		}
	}
	return queueMap
}

func (s *EVEDataServiceImpl) ensureCharacterMaps(character *model.Character) {
	if character.QualifiedPlans == nil {
		character.QualifiedPlans = make(map[string]bool)
	}
	if character.PendingPlans == nil {
		character.PendingPlans = make(map[string]bool)
	}
	if character.MissingSkills == nil {
		character.MissingSkills = make(map[string]map[string]int32)
	}
	if character.PendingFinishDates == nil {
		character.PendingFinishDates = make(map[string]*time.Time)
	}
}

type planEvaluationResult struct {
	Qualifies        bool
	Pending          bool
	MissingSkills    map[string]int32
	LatestFinishDate *time.Time
}

func (s *EVEDataServiceImpl) evaluatePlanForCharacter(
	plan model.SkillPlan,
	skillTypes map[string]model.SkillType,
	characterSkills map[int32]int32,
	skillQueueLevels map[int32]struct {
		level      int32
		finishDate *time.Time
	},
) planEvaluationResult {
	result := planEvaluationResult{
		Qualifies:     true,
		MissingSkills: make(map[string]int32),
	}

	for skillName, requiredSkill := range plan.Skills {
		skillType, exists := skillTypes[skillName]
		if !exists {
			s.logger.Warnf("Skill type not found for skill: %s", skillName)
			result.Qualifies = false
			result.MissingSkills[skillName] = int32(requiredSkill.Level)
			continue
		}

		skillID, err := strconv.Atoi(skillType.TypeID)
		if err != nil {
			s.logger.Warnf("Failed to convert TypeID to int for skill %s: %v", skillName, err)
			result.Qualifies = false
			result.MissingSkills[skillName] = int32(requiredSkill.Level)
			continue
		}

		characterLevel, hasSkill := characterSkills[int32(skillID)]
		queueInfo, inQueue := skillQueueLevels[int32(skillID)]

		requiredLevel := int32(requiredSkill.Level)

		if hasSkill && characterLevel >= requiredLevel {
			// Character already has this skill at the required level
			continue
		}

		if inQueue && queueInfo.level >= requiredLevel {
			// Skill is in queue and will meet requirement
			result.Pending = true
			result.Qualifies = false
			if result.LatestFinishDate == nil || (queueInfo.finishDate != nil && queueInfo.finishDate.After(*result.LatestFinishDate)) {
				result.LatestFinishDate = queueInfo.finishDate
			}
		} else {
			// Skill is missing or insufficient
			result.Qualifies = false
			result.MissingSkills[skillName] = requiredLevel
		}
	}

	return result
}

func (s *EVEDataServiceImpl) getStatusString(qualifies, pending bool) string {
	if qualifies {
		return "qualified"
	}
	if pending {
		return "pending"
	}
	return "missing"
}

func (s *EVEDataServiceImpl) updatePlanAndCharacterStatus(
	planStatus *model.SkillPlanWithStatus,
	character *model.Character,
	planName string,
	result planEvaluationResult,
) {
	characterStatus := model.CharacterSkillPlanStatus{
		CharacterName:     character.CharacterName,
		Status:            s.getStatusString(result.Qualifies, result.Pending),
		MissingSkills:     result.MissingSkills,
		PendingFinishDate: result.LatestFinishDate,
	}

	if result.Qualifies {
		planStatus.QualifiedCharacters = append(planStatus.QualifiedCharacters, character.CharacterName)
		character.QualifiedPlans[planName] = true
		delete(character.PendingPlans, planName)
		delete(character.MissingSkills, planName)
		delete(character.PendingFinishDates, planName)
	} else if result.Pending {
		planStatus.PendingCharacters = append(planStatus.PendingCharacters, character.CharacterName)
		character.PendingPlans[planName] = true
		character.PendingFinishDates[planName] = result.LatestFinishDate
		delete(character.QualifiedPlans, planName)
		delete(character.MissingSkills, planName)
	} else {
		planStatus.MissingCharacters = append(planStatus.MissingCharacters, character.CharacterName)
		planStatus.MissingSkills[character.CharacterName] = result.MissingSkills
		character.MissingSkills[planName] = result.MissingSkills
		delete(character.QualifiedPlans, planName)
		delete(character.PendingPlans, planName)
		delete(character.PendingFinishDates, planName)
	}

	planStatus.Characters = append(planStatus.Characters, characterStatus)
}

func (s *EVEDataServiceImpl) updateEveConversionsWithSkillNames(typeIds []int32, eveConversions map[string]string) {
	for _, typeId := range typeIds {
		skillName := s.GetSkillName(typeId)
		if skillName != "" {
			// Map skill ID to skill name (ID as string key since JSON requires string keys)
			eveConversions[strconv.Itoa(int(typeId))] = skillName
		}
	}
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

// SetAccountManagementService sets the account management service after initialization
func (s *EVEDataServiceImpl) SetAccountManagementService(accountMgmt interfaces.AccountManagementService) {
	s.accountMgmt = accountMgmt
}
