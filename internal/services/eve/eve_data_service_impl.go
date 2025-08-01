package eve

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// EVEDataServiceImpl is the full implementation without delegation
type EVEDataServiceImpl struct {
	// Core dependencies
	logger       interfaces.Logger
	httpClient   interfaces.EsiHttpClient
	authClient   interfaces.AuthClient
	accountMgmt  interfaces.AccountManagementService
	configSvc    interfaces.ConfigurationService
	
	// Repositories
	cacheRepo      interfaces.CacheRepository
	deletedRepo    interfaces.DeletedCharactersRepository
	skillRepo      interfaces.SkillRepository
	systemRepo     interfaces.SystemRepository
	eveProfileRepo interfaces.EveProfilesRepository
}

// NewEVEDataServiceImpl creates a new consolidated EVE data service implementation
func NewEVEDataServiceImpl(
	logger interfaces.Logger,
	httpClient interfaces.EsiHttpClient,
	authClient interfaces.AuthClient,
	accountMgmt interfaces.AccountManagementService,
	configSvc interfaces.ConfigurationService,
	cacheRepo interfaces.CacheRepository,
	deletedRepo interfaces.DeletedCharactersRepository,
	skillRepo interfaces.SkillRepository,
	systemRepo interfaces.SystemRepository,
	eveProfileRepo interfaces.EveProfilesRepository,
) interfaces.EVEDataService {
	svc := &EVEDataServiceImpl{
		logger:         logger,
		httpClient:     httpClient,
		authClient:     authClient,
		accountMgmt:    accountMgmt,
		configSvc:      configSvc,
		cacheRepo:      cacheRepo,
		deletedRepo:    deletedRepo,
		skillRepo:      skillRepo,
		systemRepo:     systemRepo,
		eveProfileRepo: eveProfileRepo,
	}
	return svc
}

// ESI API Operations
func (s *EVEDataServiceImpl) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	var user model.UserInfoResponse
	if err := s.httpClient.GetJSONFromURL("https://login.eveonline.com/oauth/verify", token, false, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

func (s *EVEDataServiceImpl) GetCharacter(id string) (*model.CharacterResponse, error) {
	endpoint := fmt.Sprintf("/latest/characters/%s/", id)
	resp := &model.CharacterResponse{}
	if err := s.httpClient.GetJSON(endpoint, nil, true, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *EVEDataServiceImpl) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	s.logger.Debugf("fetching skills for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skills/", characterID)
	resp := &model.CharacterSkillsResponse{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *EVEDataServiceImpl) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	s.logger.Debugf("fetching skill queue for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skillqueue/", characterID)
	var resp []model.SkillQueue
	if err := s.httpClient.GetJSON(endpoint, token, false, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *EVEDataServiceImpl) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	s.logger.Debugf("fetching location for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/location/", characterID)
	resp := &model.CharacterLocation{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return 0, err
	}

	return resp.SolarSystemID, nil
}

func (s *EVEDataServiceImpl) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	charIdToName := make(map[string]string)
	deletedChars, err := s.deletedRepo.FetchDeletedCharacters()
	if err != nil {
		s.logger.WithError(err).Info("resolve character names running without deleted characters info")
		deletedChars = []string{}
	}

	for _, id := range charIds {
		if slices.Contains(deletedChars, id) {
			continue
		}

		character, err := s.GetCharacter(id)
		if err != nil {
			s.logger.Warnf("failed to retrieve name for %s", id)
			var customErr *flyErrors.CustomError
			if errors.As(err, &customErr) && customErr.StatusCode == http.StatusNotFound {
				s.logger.Warnf("adding %s to deleted characters", id)
				deletedChars = append(deletedChars, id)
			}
		} else {
			charIdToName[id] = character.Name
		}
	}

	if saveErr := s.deletedRepo.SaveDeletedCharacters(deletedChars); saveErr != nil {
		s.logger.Warnf("failed to save deleted characters %v", saveErr)
	}
	if err := s.SaveCache(); err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdToName, nil
}

func (s *EVEDataServiceImpl) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	endpoint := fmt.Sprintf("/latest/corporations/%d/", id)
	resp := &model.Corporation{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *EVEDataServiceImpl) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	endpoint := fmt.Sprintf("/latest/alliances/%d/", id)
	resp := &model.Alliance{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Character Management
func (s *EVEDataServiceImpl) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	s.logger.Debugf("Processing identity for character ID: %d", charIdentity.Character.CharacterID)

	user, err := s.GetUserInfo(&charIdentity.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	s.logger.Debugf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	characterResponse, err := s.GetCharacter(strconv.FormatInt(charIdentity.Character.CharacterID, 10))
	if err != nil {
		s.logger.Warnf("Failed to get character %s: %v", charIdentity.Character.CharacterName, err)
	}

	skills, err := s.GetCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Warnf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}
	s.logger.Debugf("Fetched %d skills for character %d", len(skills.Skills), charIdentity.Character.CharacterID)

	skillQueue, err := s.GetCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Warnf("Failed to get eve queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	s.logger.Debugf("Fetched %d eve queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	characterLocation, err := s.GetCharacterLocation(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Warnf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	}

	corporationName := ""
	allianceName := ""
	if characterResponse != nil {
		characterCorporation, err := s.GetCorporation(int64(characterResponse.CorporationID), &charIdentity.Token)
		if err != nil {
			s.logger.Warnf("Failed to get corporation for corporation %d: %v", characterResponse.CorporationID, err)
		} else {
			corporationName = characterCorporation.Name
		}
		if characterCorporation != nil && characterCorporation.AllianceID != 0 {
			characterAlliance, err := s.GetAlliance(int64(characterCorporation.AllianceID), &charIdentity.Token)
			if err != nil {
				s.logger.Warnf("Failed to get alliance for character %s: %v", characterCorporation.AllianceID, err)
			} else {
				allianceName = characterAlliance.Name
			}
		}
	}

	s.logger.Debugf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	// Update charIdentity with fetched data
	s.logger.Debugf("updating %s", user.CharacterName)
	charIdentity.Character.UserInfoResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = s.systemRepo.GetSystemName(charIdentity.Character.Location)
	charIdentity.MCT = s.isCharacterTraining(*skillQueue)
	if charIdentity.MCT {
		charIdentity.Training = s.GetSkillName(charIdentity.Character.SkillQueue[0].SkillID)
	}
	charIdentity.CorporationName = corporationName
	charIdentity.AllianceName = allianceName

	// Initialize maps if nil
	if charIdentity.Character.QualifiedPlans == nil {
		charIdentity.Character.QualifiedPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingPlans == nil {
		charIdentity.Character.PendingPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingFinishDates == nil {
		charIdentity.Character.PendingFinishDates = make(map[string]*time.Time)
	}
	if charIdentity.Character.MissingSkills == nil {
		charIdentity.Character.MissingSkills = make(map[string]map[string]int32)
	}

	err = s.SaveCache()
	if err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdentity, nil
}

func (s *EVEDataServiceImpl) isCharacterTraining(queue []model.SkillQueue) bool {
	for _, q := range queue {
		if q.StartDate != nil && q.FinishDate != nil && q.FinishDate.After(time.Now()) {
			s.logger.Debugf("training - start %s, finish %s, eve %d", q.StartDate, q.FinishDate, q.SkillID)
			return true
		}
	}
	return false
}

func (s *EVEDataServiceImpl) DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error) {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return false, nil, err
	}

	for _, account := range accounts {
		for _, char := range account.Characters {
			if char.Character.CharacterID == characterID {
				return true, &char, nil
			}
		}
	}
	return false, nil, nil
}

func (s *EVEDataServiceImpl) UpdateCharacterFields(characterID int64, updates map[string]interface{}) error {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to fetch accounts: %w", err)
	}

	// Find the character
	var charIdentity *model.CharacterIdentity
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				charIdentity = &accounts[i].Characters[j]
				break
			}
		}
		if charIdentity != nil {
			break
		}
	}
	
	if charIdentity == nil {
		return fmt.Errorf("character not found")
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "Role":
			roleStr, ok := value.(string)
			if !ok {
				return fmt.Errorf("role must be a string")
			}
			// Update roles via configService
			if err := s.configSvc.UpdateRoles(roleStr); err != nil {
				s.logger.Infof("Failed to update roles: %v", err)
			}
			charIdentity.Role = roleStr

		case "MCT":
			mctBool, ok := value.(bool)
			if !ok {
				return fmt.Errorf("MCT must be boolean")
			}
			charIdentity.MCT = mctBool

		default:
			return fmt.Errorf("unknown update field: %s", key)
		}
	}

	if err := s.accountMgmt.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (s *EVEDataServiceImpl) RemoveCharacter(characterID int64) error {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return err
	}

	for i := range accounts {
		var updatedChars []model.CharacterIdentity
		for _, char := range accounts[i].Characters {
			if char.Character.CharacterID != characterID {
				updatedChars = append(updatedChars, char)
			}
		}
		if len(updatedChars) != len(accounts[i].Characters) {
			accounts[i].Characters = updatedChars
			return s.accountMgmt.SaveAccounts(accounts)
		}
	}
	
	return fmt.Errorf("character not found")
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
		CharacterName: character.CharacterName,
		Status:        s.getStatusString(result.Qualifies, result.Pending),
		MissingSkills: result.MissingSkills,
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
			eveConversions[skillName] = strconv.Itoa(int(typeId))
		}
	}
}

// Eve Profile Management
func (s *EVEDataServiceImpl) LoadCharacterSettings() ([]model.EveProfile, error) {
	settingsDir, err := s.configSvc.GetSettingsDir()
	if err != nil {
		return nil, err
	}

	subDirs, err := s.eveProfileRepo.GetSubDirectories(settingsDir)
	if err != nil {
		return nil, err
	}

	var settingsData []model.EveProfile
	allCharIDs := make(map[string]struct{})

	for _, sd := range subDirs {
		rawFiles, err := s.eveProfileRepo.ListSettingsFiles(sd, settingsDir)
		if err != nil {
			s.logger.Warnf("Error fetching settings files for subDir %s: %v", sd, err)
			continue
		}

		var charFiles []model.CharFile
		var userFiles []model.UserFile

		for _, rf := range rawFiles {
			if rf.IsChar {
				// Just record charId for later ESI resolution
				allCharIDs[rf.CharOrUserID] = struct{}{}
				charFiles = append(charFiles, model.CharFile{
					File:   rf.FileName,
					CharId: rf.CharOrUserID,
					Name:   "CharID:" + rf.CharOrUserID, // Temporary name, will update after ESI lookup
					Mtime:  rf.Mtime,
				})
			} else {
				friendlyName := rf.CharOrUserID
				if savedName, ok := s.accountMgmt.GetAccountNameByID(rf.CharOrUserID); ok {
					friendlyName = savedName
				}
				userFiles = append(userFiles, model.UserFile{
					File:   rf.FileName,
					UserId: rf.CharOrUserID,
					Name:   friendlyName,
					Mtime:  rf.Mtime,
				})
			}
		}

		settingsData = append(settingsData, model.EveProfile{
			Profile:            sd,
			AvailableCharFiles: charFiles,
			AvailableUserFiles: userFiles,
		})
	}

	// Resolve character names via ESI
	var charIdList []string
	for id := range allCharIDs {
		charIdList = append(charIdList, id)
	}
	charIdToName, err := s.ResolveCharacterNames(charIdList)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve character names: %w", err)
	}

	// Update character files with resolved names
	for si, sd := range settingsData {
		var filteredChars []model.CharFile
		for _, cf := range sd.AvailableCharFiles {
			if name, ok := charIdToName[cf.CharId]; ok && name != "" {
				cf.Name = name
				filteredChars = append(filteredChars, cf)
			}
		}
		settingsData[si].AvailableCharFiles = filteredChars
	}

	return settingsData, nil
}

func (s *EVEDataServiceImpl) BackupDir(targetDir, backupDir string) error {
	// Backup the EVE "settings_" directories
	if err := s.eveProfileRepo.BackupDirectory(targetDir, backupDir); err != nil {
		return err
	}
	
	// Also backup JSON files from config
	return s.configSvc.BackupJSONFiles(backupDir)
}

func (s *EVEDataServiceImpl) SyncDir(subDir, charId, userId string) (int, int, error) {
	settingsDir, err := s.configSvc.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}

	return s.eveProfileRepo.SyncSubdirectory(subDir, userId, charId, settingsDir)
}

func (s *EVEDataServiceImpl) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	settingsDir, err := s.configSvc.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}
	if settingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return s.eveProfileRepo.SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir)
}

// Cache Management
func (s *EVEDataServiceImpl) SaveCache() error {
	return s.cacheRepo.SaveApiCache()
}

func (s *EVEDataServiceImpl) LoadCache() error {
	return s.cacheRepo.LoadApiCache()
}

// SaveEsiCache saves the ESI cache (alias for SaveCache to implement ESIService interface)
func (s *EVEDataServiceImpl) SaveEsiCache() error {
	return s.SaveCache()
}