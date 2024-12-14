package eve

import (
	"bufio"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"strconv"
	"strings"
	"time"
)

var _ interfaces.SkillService = (*skillService)(nil)

type skillService struct {
	logger    interfaces.Logger
	skillRepo interfaces.SkillRepository
}

// NewSkillService  returns a new SettingsService with a config
func NewSkillService(logger interfaces.Logger, skillRepo interfaces.SkillRepository) interfaces.SkillService {
	return &skillService{logger: logger, skillRepo: skillRepo}
}

var romanToInt = map[string]int{
	"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5,
}

func (s *skillService) GetSkillPlanFile(planName string) ([]byte, error) {
	return s.skillRepo.GetSkillPlanFile(planName)
}

func (s *skillService) DeleteSkillPlan(name string) error {
	return s.skillRepo.DeleteSkillPlan(name)
}

func (s *skillService) ParseAndSaveSkillPlan(contents, name string) error {
	skills := s.parseSkillPlanContents(contents)
	return s.skillRepo.SaveSkillPlan(name, skills)
}

func (s *skillService) CheckIfDuplicatePlan(name string) bool {
	plans := s.skillRepo.GetSkillPlans()
	for _, plan := range plans {
		if plan.Name == name {
			return true
		}
	}
	return false
}

// parseSkillPlanContents takes the contents as a string and parses it into a map of skills.
func (s *skillService) parseSkillPlanContents(contents string) map[string]model.Skill {
	skills := make(map[string]model.Skill)
	scanner := bufio.NewScanner(strings.NewReader(contents))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue // Skip empty lines
		}

		// Find the last whitespace to separate eve name from eve level
		lastSpaceIndex := strings.LastIndex(line, " ")
		if lastSpaceIndex == -1 {
			continue // Skip lines that don't have a space
		}

		// Separate eve name and level
		skillName := line[:lastSpaceIndex]
		skillLevelStr := line[lastSpaceIndex+1:]

		// Parse eve level, handling Roman numerals if necessary
		skillLevel, err := parseSkillLevel(skillLevelStr)
		if err != nil {
			s.logger.Warnf("Invalid eve level '%s'; skipping line.\n", skillLevelStr)
			continue // Skip lines with invalid levels
		}

		// Check if the eve already exists and add/update if necessary
		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}

	return skills
}

// parseSkillLevel converts either a Roman numeral or integer string to an integer.
func parseSkillLevel(levelStr string) (int, error) {
	if val, ok := romanToInt[levelStr]; ok {
		return val, nil
	}
	return strconv.Atoi(levelStr) // Fall back to numeric conversion
}

func (s *skillService) GetPlanAndConversionData(
	accounts []model.Account,
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType,
) (map[string]model.SkillPlanWithStatus, map[string]string) {

	// Step 1: Initialize updatedSkillPlans and eveConversions
	updatedSkillPlans := s.initializeUpdatedPlans(skillPlans)
	eveConversions := s.initializeEveConversions(skillPlans, skillTypes)

	// Step 2: Process all accounts and characters
	typeIds := s.processAccountsAndCharacters(accounts, skillPlans, skillTypes, updatedSkillPlans)

	// Step 3: Convert skill IDs into names and update eveConversions
	s.updateEveConversionsWithSkillNames(typeIds, eveConversions)

	return updatedSkillPlans, eveConversions
}

func (s *skillService) initializeUpdatedPlans(skillPlans map[string]model.SkillPlan) map[string]model.SkillPlanWithStatus {
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

func (s *skillService) initializeEveConversions(
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

func (s *skillService) processAccountsAndCharacters(
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

			s.ensureCharacterMaps(character)

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

func (s *skillService) mapCharacterSkills(character model.Character, typeIds *[]int32) map[int32]int32 {
	skillsMap := make(map[int32]int32)
	for _, skill := range character.Skills {
		skillsMap[skill.SkillID] = skill.TrainedSkillLevel
		*typeIds = append(*typeIds, skill.SkillID)
	}
	return skillsMap
}

func (s *skillService) mapSkillQueueLevels(character model.Character) map[int32]struct {
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

func (s *skillService) ensureCharacterMaps(character model.Character) {
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

func (s *skillService) evaluatePlanForCharacter(
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
			s.logger.Errorf("Error: Skill '%s' does not exist in eve types", skillName)
			result.Qualifies = false
			continue
		}

		skillID, err := strconv.Atoi(skillType.TypeID)
		if err != nil {
			s.logger.Errorf("Error: Converting eve type ID '%s' for eve '%s': %v", skillType.TypeID, skillName, err)
			result.Qualifies = false
			continue
		}

		requiredLevel := int32(requiredSkill.Level)
		characterLevel, hasSkill := characterSkills[int32(skillID)]
		queued, inQueue := skillQueueLevels[int32(skillID)]

		switch {
		case hasSkill && characterLevel >= requiredLevel:
			// Already qualified for this skill
		case inQueue && queued.level >= requiredLevel:
			// Pending this skill
			result.Pending = true
			if result.LatestFinishDate == nil || (queued.finishDate != nil && queued.finishDate.After(*result.LatestFinishDate)) {
				result.LatestFinishDate = queued.finishDate
			}
		default:
			// Missing this skill
			result.Qualifies = false
			result.MissingSkills[skillName] = requiredLevel
		}
	}

	return result
}
func (s *skillService) updatePlanAndCharacterStatus(
	plan *model.SkillPlanWithStatus,
	character *model.Character,
	planName string,
	res planEvaluationResult,
) {
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

	characterSkillStatus := model.CharacterSkillPlanStatus{
		CharacterName:     character.CharacterName,
		Status:            getStatus(res.Qualifies, res.Pending),
		MissingSkills:     res.MissingSkills,
		PendingFinishDate: res.LatestFinishDate,
	}

	if res.Qualifies && !res.Pending {
		plan.QualifiedCharacters = append(plan.QualifiedCharacters, character.CharacterName)
		character.QualifiedPlans[planName] = true
	}

	if res.Pending {
		plan.PendingCharacters = append(plan.PendingCharacters, character.CharacterName)
		character.PendingPlans[planName] = true
		character.PendingFinishDates[planName] = res.LatestFinishDate
	}

	if len(res.MissingSkills) > 0 {
		plan.MissingSkills[character.CharacterName] = res.MissingSkills
		character.MissingSkills[planName] = res.MissingSkills
	}

	plan.Characters = append(plan.Characters, characterSkillStatus)
}

func (s *skillService) updateEveConversionsWithSkillNames(typeIds []int32, eveConversions map[string]string) {
	for _, skillId := range typeIds {
		name := s.GetSkillName(skillId)
		if name != "" {
			eveConversions[strconv.FormatInt(int64(skillId), 10)] = name
		}
	}
}

func getStatus(qualifies bool, pending bool) string {
	if qualifies && !pending {
		return "Qualified"
	} else if pending {
		return "Pending"
	}
	return "Not Qualified"
}

func (s *skillService) GetSkillPlans() map[string]model.SkillPlan {
	return s.skillRepo.GetSkillPlans()
}

func (s *skillService) GetSkillTypes() map[string]model.SkillType {
	return s.skillRepo.GetSkillTypes()
}

func (s *skillService) GetSkillTypeByID(id string) (model.SkillType, bool) {
	return s.skillRepo.GetSkillTypeByID(id)
}

func (s *skillService) GetSkillName(skillID int32) string {
	skill, ok := s.skillRepo.GetSkillTypeByID(strconv.FormatInt(int64(skillID), 10))
	if !ok {
		s.logger.Warnf("Skill ID %d not found", skillID)
		return ""
	}
	return skill.TypeName
}
