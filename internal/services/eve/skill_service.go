package eve

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

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

func (s *skillService) GetMatchingSkillPlans(
	accounts []model.Account,
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType,
) map[string]model.SkillPlanWithStatus {

	updatedSkillPlans := make(map[string]model.SkillPlanWithStatus)

	// Initialize updatedSkillPlans with empty QualifiedCharacters, PendingCharacters, and Characters
	for planName, plan := range skillPlans {
		updatedSkillPlans[planName] = model.SkillPlanWithStatus{
			Name:                plan.Name,
			Skills:              plan.Skills,
			QualifiedCharacters: []string{},
			PendingCharacters:   []string{},
			MissingSkills:       make(map[string]map[string]int32), // Per character missing skills
			Characters:          []model.CharacterSkillPlanStatus{},
		}
	}

	// Process each account and its characters
	for _, account := range accounts {
		for _, characterData := range account.Characters {
			character := characterData.Character
			characterSkills := make(map[int32]int32)

			// Map character's current skills for lookup
			for _, skill := range character.Skills {
				characterSkills[skill.SkillID] = skill.TrainedSkillLevel
			}

			// Track eve queue levels for pending skills
			skillQueueLevels := make(map[int32]struct {
				level      int32
				finishDate *time.Time
			})
			for _, queuedSkill := range character.SkillQueue {
				if current, exists := skillQueueLevels[queuedSkill.SkillID]; !exists || queuedSkill.FinishedLevel > current.level {
					skillQueueLevels[queuedSkill.SkillID] = struct {
						level      int32
						finishDate *time.Time
					}{level: queuedSkill.FinishedLevel, finishDate: queuedSkill.FinishDate}
				}
			}

			// Initialize QualifiedPlans, PendingPlans, and MissingSkills for the character
			if character.QualifiedPlans == nil {
				character.QualifiedPlans = make(map[string]bool)
			}
			if character.PendingPlans == nil {
				character.PendingPlans = make(map[string]bool)
			}
			if character.MissingSkills == nil {
				character.MissingSkills = make(map[string]map[string]int32)
			}

			// Check matching eve plans for the current character
			for planName, plan := range skillPlans {
				qualifies := true
				pending := false
				missingSkills := make(map[string]int32)
				var latestFinishDate *time.Time

				for skillName, requiredSkill := range plan.Skills {
					// Map skillName to skillID via skillTypes
					skillType, exists := skillTypes[skillName]
					if !exists {
						s.logger.Errorf("Error: Skill '%s' does not exist in eve types", skillName)
						qualifies = false
						continue
					}

					skillID, err := strconv.Atoi(skillType.TypeID)
					if err != nil {
						s.logger.Errorf("Error: Converting eve type ID '%s' for eve '%s': %v", skillType.TypeID, skillName, err)
						qualifies = false
						continue
					}

					requiredLevel := int32(requiredSkill.Level)
					characterLevel, hasSkill := characterSkills[int32(skillID)]
					queued, inQueue := skillQueueLevels[int32(skillID)]

					// Compare current eve and eve queue for qualification
					if hasSkill && characterLevel >= requiredLevel {
						continue // Qualified for this eve
					} else if inQueue && queued.level >= requiredLevel {
						pending = true
						if latestFinishDate == nil || (queued.finishDate != nil && queued.finishDate.After(*latestFinishDate)) {
							latestFinishDate = queued.finishDate
						}
					} else {
						qualifies = false
						missingSkills[skillName] = requiredLevel
					}
				}

				// Add this character's status to the eve plan's Characters list
				characterSkillStatus := model.CharacterSkillPlanStatus{
					CharacterName:     character.CharacterName,
					Status:            getStatus(qualifies, pending),
					MissingSkills:     missingSkills, // Store missing skills correctly
					PendingFinishDate: latestFinishDate,
				}

				modifiedPlan := updatedSkillPlans[planName]

				// If qualifies and not pending, add to QualifiedCharacters
				if qualifies && !pending {
					modifiedPlan.QualifiedCharacters = append(modifiedPlan.QualifiedCharacters, character.CharacterName)
					character.QualifiedPlans[planName] = true // Update the character's QualifiedPlans
				}

				// If pending, add to PendingCharacters
				if pending {
					modifiedPlan.PendingCharacters = append(modifiedPlan.PendingCharacters, character.CharacterName)
					characterSkillStatus.PendingFinishDate = latestFinishDate
					character.PendingPlans[planName] = true // Update the character's PendingPlans
					character.PendingFinishDates[planName] = latestFinishDate
				}

				// If there are missing skills, store them in the MissingSkills map
				if len(missingSkills) > 0 {
					if modifiedPlan.MissingSkills == nil {
						modifiedPlan.MissingSkills = make(map[string]map[string]int32)
					}
					modifiedPlan.MissingSkills[character.CharacterName] = missingSkills
					character.MissingSkills[planName] = missingSkills // Update MissingSkills for the character
				}

				// Add the character's eve status to the plan's character list
				modifiedPlan.Characters = append(modifiedPlan.Characters, characterSkillStatus)

				// Update the eve plan with the new status
				updatedSkillPlans[planName] = modifiedPlan
			}
		}
	}

	return updatedSkillPlans
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

func (sk *skillService) GetSkillName(skillID int32) string {
	skill, ok := sk.skillRepo.GetSkillTypeByID(strconv.FormatInt(int64(skillID), 10))
	if !ok {
		sk.logger.Warnf("Skill ID %d not found", skillID)
		return ""
	}
	return skill.TypeName
}
