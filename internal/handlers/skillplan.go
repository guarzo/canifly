package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/service/skillplan"
	"github.com/guarzo/canifly/internal/utils/xlog"

	"github.com/guarzo/canifly/internal/model"
)

var romanToInt = map[string]int{
	"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5,
}

// SkillPlanFileHandler serves skill plan file contents based on the provided plan name.
func SkillPlanFileHandler(w http.ResponseWriter, r *http.Request) {
	planName := r.URL.Query().Get("planName")
	if planName == "" {
		http.Error(w, "Missing planName parameter", http.StatusBadRequest)
		return
	}

	// Append .txt extension to the file name
	planName += ".txt"
	xlog.Logf("Attempting to serve skill plan file: %s", planName)

	// Get the writable plans directory path
	skillPlanDir, err := skillplan.GetWritablePlansPath()
	if err != nil {
		xlog.Logf("Failed to retrieve skill plan directory: %v", err)
		http.Error(w, fmt.Sprintf("Failed to retrieve skill plan directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct the full file path
	filePath := filepath.Join(skillPlanDir, planName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Skill plan file not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to read skill plan file: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(content)
	if writeErr != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", writeErr), http.StatusInternalServerError)
	}
}

func SaveSkillPlanHandler(w http.ResponseWriter, r *http.Request) {
	// Define struct with contents as a string
	var requestData struct {
		PlanName string `json:"name"`
		Contents string `json:"contents"`
	}

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		xlog.Logf("Failed to parse JSON body: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	planName := requestData.PlanName
	contents := requestData.Contents

	if planName == "" {
		xlog.Log("planName parameter missing")
		http.Error(w, "Missing planName parameter", http.StatusBadRequest)
		return
	}

	// Convert contents from string to the map format expected by the skill plan.
	skills := parseSkillPlanContents(contents)

	// Save the skill plan
	if err := skillplan.SaveSkillPlan(planName, skills); err != nil {
		xlog.Logf("Failed to save skill plan: %v", err)
		http.Error(w, "Failed to save skill plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// parseSkillPlanContents takes the contents as a string and parses it into a map of skills.
func parseSkillPlanContents(contents string) map[string]model.Skill {
	skills := make(map[string]model.Skill)
	scanner := bufio.NewScanner(strings.NewReader(contents))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue // Skip empty lines
		}

		// Find the last whitespace to separate skill name from skill level
		lastSpaceIndex := strings.LastIndex(line, " ")
		if lastSpaceIndex == -1 {
			continue // Skip lines that don't have a space
		}

		// Separate skill name and level
		skillName := line[:lastSpaceIndex]
		skillLevelStr := line[lastSpaceIndex+1:]

		// Parse skill level, handling Roman numerals if necessary
		skillLevel, err := parseSkillLevel(skillLevelStr)
		if err != nil {
			xlog.Logf("Invalid skill level '%s'; skipping line.\n", skillLevelStr)
			continue // Skip lines with invalid levels
		}

		// Check if the skill already exists and add/update if necessary
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

func getMatchingSkillPlans(
	characters map[int64]model.CharacterIdentity,
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType, // Retain the skillTypes map parameter
) (map[int64]model.CharacterIdentity, map[string]model.SkillPlan) {

	updatedCharacters := make(map[int64]model.CharacterIdentity)
	updatedSkillPlans := make(map[string]model.SkillPlan)

	// Initialize updatedSkillPlans with empty QualifiedCharacters and PendingCharacters
	for planName, plan := range skillPlans {
		plan.QualifiedCharacters = []string{}
		plan.PendingCharacters = []string{}
		updatedSkillPlans[planName] = plan
	}

	for charID, characterData := range characters {
		character := characterData.Character
		character.QualifiedPlans = make(map[string]bool)
		character.PendingPlans = make(map[string]bool)
		character.PendingFinishDates = make(map[string]*time.Time)
		character.MissingSkills = make(map[string]map[string]int32)

		// Map character's current skills and skill queue for lookup
		characterSkills := make(map[int32]int32)
		for _, skill := range character.Skills {
			characterSkills[skill.SkillID] = skill.TrainedSkillLevel
		}

		// Track maximum levels of skills in the queue and their finish dates
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

		for planName, plan := range updatedSkillPlans {
			qualifies := true
			pending := false
			missingSkills := make(map[string]int32)
			var latestFinishDate *time.Time

			for skillName, requiredSkill := range plan.Skills {

				// Map skillName to skillID via skillTypes
				skillType, exists := skillTypes[skillName]
				if !exists {
					xlog.Logf("Error: Skill '%s' does not exist in skill types", skillName)
					qualifies = false
					continue // Instead of break, use continue to move to the next skill
				}

				skillID, err := strconv.Atoi(skillType.TypeID)
				if err != nil {
					xlog.Logf("Error: Converting skill type ID '%s' for skill '%s': %v", skillType.TypeID, skillName, err)
					qualifies = false
					continue
				}

				requiredLevel := int32(requiredSkill.Level)
				characterLevel, hasSkill := characterSkills[int32(skillID)]
				queued, inQueue := skillQueueLevels[int32(skillID)]

				// Compare current skill and skill queue for qualification
				if hasSkill && characterLevel >= requiredLevel {
					continue
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

			// Update Qualified and Pending status
			if qualifies && !pending {
				character.QualifiedPlans[planName] = true
				modifiedPlan := updatedSkillPlans[planName]
				modifiedPlan.QualifiedCharacters = append(modifiedPlan.QualifiedCharacters, character.CharacterName)
				updatedSkillPlans[planName] = modifiedPlan
			} else if pending {
				character.PendingPlans[planName] = true
				modifiedPlan := updatedSkillPlans[planName]
				modifiedPlan.PendingCharacters = append(modifiedPlan.PendingCharacters, character.CharacterName)
				updatedSkillPlans[planName] = modifiedPlan
				character.PendingFinishDates[planName] = latestFinishDate
			}

			// If there are any missing skills, add them to the character's MissingSkills map
			if len(missingSkills) > 0 {
				character.MissingSkills[planName] = missingSkills
			}
		}
		characterData.Character = character
		updatedCharacters[charID] = characterData
	}

	return updatedCharacters, updatedSkillPlans
}

// DeleteSkillPlanHandler handles the deletion of a skill plan.
func DeleteSkillPlanHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a DELETE request
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the planName from the URL query parameters
	planName := r.URL.Query().Get("planName")
	if planName == "" {
		xlog.Log("planName parameter missing")
		http.Error(w, "Missing planName parameter", http.StatusBadRequest)
		return
	}

	// Delete the skill plan
	if err := skillplan.DeleteSkillPlan(planName); err != nil {
		xlog.Logf("Failed to delete skill plan: %v", err)
		http.Error(w, "Failed to delete skill plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
