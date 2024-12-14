// model/eve.go

package model

import (
	"time"
)

type EveData struct {
	EveProfiles    []EveProfile
	SkillPlans     map[string]SkillPlanWithStatus
	EveConversions map[string]string // converts skill id to skill name
}

// EveProfile is the data from the eve settings
type EveProfile struct {
	Profile            string     `json:"profile"`            // eve profile name
	AvailableCharFiles []CharFile `json:"availableCharFiles"` // character files for a given profile
	AvailableUserFiles []UserFile `json:"availableUserFiles"` // user files for a given profile
}

// RawFileInfo represents basic information extracted from an EVE settings file.
type RawFileInfo struct {
	FileName     string
	CharOrUserID string
	IsChar       bool // true if char file, false if user file
	Mtime        string
}

type Station struct {
	SystemID int64  `json:"system_id"`
	ID       int64  `json:"station_id"`
	Name     string `json:"station_name"`
}

type Structure struct {
	Name     string `json:"name"`
	OwnerID  int64  `json:"owner_id"`
	SystemID int64  `json:"solar_system_id"`
	TypeID   int64  `json:"type_id"`
}

type CharacterLocation struct {
	SolarSystemID int64 `json:"solar_system_id"`
	StructureID   int64 `json:"structure_id"`
}

type CloneLocation struct {
	HomeLocation struct {
		LocationID   int64  `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"home_location"`
	JumpClones []struct {
		Implants     []int  `json:"implants"`
		JumpCloneID  int64  `json:"jump_clone_id"`
		LocationID   int64  `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"jump_clones"`
}

// SkillPlanWithStatus holds detailed information about each eve plan
type SkillPlanWithStatus struct {
	Name                string
	TypeId              int64 // used for image lookup
	Skills              map[string]Skill
	QualifiedCharacters []string
	PendingCharacters   []string
	MissingSkills       map[string]map[string]int32 // Missing skills by character
	Characters          []CharacterSkillPlanStatus  // List of characters with their status for this eve plan
}

// CharacterSkillPlanStatus represents a character's status for a specific eve plan
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

// Skill represents a eve with a name and level.
type Skill struct {
	Name  string `json:"Name"`
	Level int    `json:"Level"`
}

// SkillPlan represents a eve plan, with the plan name and a map of unique skills.
type SkillPlan struct {
	Name                string           `json:"Name"`
	Skills              map[string]Skill `json:"Skills"`
	QualifiedCharacters []string         `json:"QualifiedCharacters"`
	PendingCharacters   []string         `json:"PendingCharacters"`
}

// SkillType represents a eve with typeID, typeName, and description.
type SkillType struct {
	TypeID      string
	TypeName    string
	Description string
}

type CharFile struct {
	File   string `json:"file"`
	CharId string `json:"charId"`
	Name   string `json:"name"`
	Mtime  string `json:"mtime"`
}

type UserFile struct {
	File   string `json:"file"`
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Mtime  string `json:"mtime"`
}
