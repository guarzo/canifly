package skillplan

import (
	"strconv"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// Compile-time interface check.
var _ interfaces.SkillPlanService = (*Service)(nil)

// Service implements interfaces.SkillPlanService.
type Service struct {
	logger    interfaces.Logger
	skillRepo interfaces.SkillRepository
}

// NewService constructs a SkillPlanService.
func NewService(logger interfaces.Logger, skillRepo interfaces.SkillRepository) *Service {
	return &Service{logger: logger, skillRepo: skillRepo}
}

// Skill Plan Management
func (s *Service) GetSkillPlans() map[string]model.SkillPlan {
	return s.skillRepo.GetSkillPlans()
}

func (s *Service) GetSkillName(id int32) string {
	skillType, found := s.skillRepo.GetSkillTypeByID(strconv.Itoa(int(id)))
	if found {
		return skillType.TypeName
	}
	return ""
}

func (s *Service) GetSkillTypes() map[string]model.SkillType {
	return s.skillRepo.GetSkillTypes()
}

func (s *Service) CheckIfDuplicatePlan(name string) bool {
	plans := s.skillRepo.GetSkillPlans()
	_, exists := plans[name]
	return exists
}

func (s *Service) ParseAndSaveSkillPlan(contents, name string) error {
	skills := s.parseSkillPlan(contents)
	return s.skillRepo.SaveSkillPlan(name, skills)
}

func (s *Service) parseSkillPlan(contents string) map[string]model.Skill {
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

func (s *Service) GetSkillPlanFile(name string) ([]byte, error) {
	return s.skillRepo.GetSkillPlanFile(name)
}

func (s *Service) DeleteSkillPlan(name string) error {
	return s.skillRepo.DeleteSkillPlan(name)
}

func (s *Service) ListSkillPlans() ([]string, error) {
	plans := s.GetSkillPlans()
	planNames := make([]string, 0, len(plans))
	for name := range plans {
		planNames = append(planNames, name)
	}
	return planNames, nil
}

func (s *Service) RefreshRemotePlans() error {
	// Reload skill plans from repository (which will trigger GitHub download if configured)
	return s.skillRepo.LoadSkillPlans()
}

func (s *Service) GetSkillTypeByID(id string) (model.SkillType, bool) {
	return s.skillRepo.GetSkillTypeByID(id)
}

func (s *Service) GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string) {
	// Step 1: Initialize updatedSkillPlans and eveConversions
	updatedSkillPlans := s.initializeUpdatedPlans(skillPlans)
	eveConversions := s.initializeEveConversions(skillPlans, skillTypes)

	// Step 2: Process all accounts and characters
	typeIds := s.processAccountsAndCharacters(accounts, skillPlans, skillTypes, updatedSkillPlans)

	// Step 3: Convert skill IDs into names and update eveConversions
	s.updateEveConversionsWithSkillNames(typeIds, eveConversions)

	return updatedSkillPlans, eveConversions
}

func (s *Service) initializeUpdatedPlans(skillPlans map[string]model.SkillPlan) map[string]model.SkillPlanWithStatus {
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

func (s *Service) initializeEveConversions(
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

func (s *Service) processAccountsAndCharacters(
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

func (s *Service) mapCharacterSkills(character model.Character, typeIds *[]int32) map[int32]int32 {
	skillsMap := make(map[int32]int32)
	for _, skill := range character.Skills {
		skillsMap[skill.SkillID] = skill.TrainedSkillLevel
		*typeIds = append(*typeIds, skill.SkillID)
	}
	return skillsMap
}

func (s *Service) mapSkillQueueLevels(character model.Character) map[int32]struct {
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

func (s *Service) ensureCharacterMaps(character *model.Character) {
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

func (s *Service) evaluatePlanForCharacter(
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

func (s *Service) getStatusString(qualifies, pending bool) string {
	if qualifies {
		return "qualified"
	}
	if pending {
		return "pending"
	}
	return "missing"
}

func (s *Service) updatePlanAndCharacterStatus(
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

func (s *Service) updateEveConversionsWithSkillNames(typeIds []int32, eveConversions map[string]string) {
	for _, typeId := range typeIds {
		skillName := s.GetSkillName(typeId)
		if skillName != "" {
			// Map skill ID to skill name (ID as string key since JSON requires string keys)
			eveConversions[strconv.Itoa(int(typeId))] = skillName
		}
	}
}
