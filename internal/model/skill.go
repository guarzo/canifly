// model/skill.go
package model

import (
	"time"
)

// SkillPlanWithStatus holds detailed information about each skill plan
type SkillPlanWithStatus struct {
	Name                string
	Skills              map[string]Skill
	QualifiedCharacters []string
	PendingCharacters   []string
	MissingSkills       map[string]map[string]int32 // Missing skills by character
	Characters          []CharacterSkillPlanStatus  // List of characters with their status for this skill plan
}

// CharacterSkillPlanStatus represents a character's status for a specific skill plan
type CharacterSkillPlanStatus struct {
	CharacterName     string
	Status            string // "qualified", "pending", "missing"
	MissingSkills     map[string]int32
	PendingFinishDate *time.Time
}

type SkillResponse struct {
	ActiveSkillLevel   int32 `json:"active_skill_level"`
	SkillID            int32 `json:"skill_id"`
	SkillpointsInSkill int64 `json:"skillpoints_in_skill"`
	TrainedSkillLevel  int32 `json:"trained_skill_level"`
}

type SkillQueue struct {
	FinishDate      *time.Time `json:"finish_date,omitempty"`
	FinishedLevel   int32      `json:"finished_level"`
	LevelEndSP      int32      `json:"level_end_sp"`
	LevelStartSP    int32      `json:"level_start_sp"`
	QueuePosition   int32      `json:"queue_position"`
	SkillID         int32      `json:"skill_id"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	TrainingStartSP int32      `json:"training_start_sp"`
}

// Skill represents a skill with a name and level.
type Skill struct {
	Name  string `json:"Name"`
	Level int    `json:"Level"`
}

// SkillPlan represents a skill plan, with the plan name and a map of unique skills.
type SkillPlan struct {
	Name                string           `json:"Name"`
	Skills              map[string]Skill `json:"Skills"`
	QualifiedCharacters []string         `json:"QualifiedCharacters"`
	PendingCharacters   []string         `json:"PendingCharacters"`
}

// SkillType represents a skill with typeID, typeName, and description.
type SkillType struct {
	TypeID      string
	TypeName    string
	Description string
}
