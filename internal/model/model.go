package model

import (
	"time"

	"golang.org/x/oauth2"
)

// HomeData is the structure used to pass the data to the frontend
type HomeData struct {
	Title        string
	LoggedIn     bool
	Identities   map[int64]CharacterIdentity
	SkillPlans   map[string]SkillPlanWithStatus // Skill plans with character status
	MainIdentity int64
}

// SkillPlanWithStatus holds detailed information about each skill plan
type SkillPlanWithStatus struct {
	Name                string
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

// CharacterIdentity structure contains nested Character struct
type CharacterIdentity struct {
	Token     oauth2.Token
	Character Character
}

type Character struct {
	BaseCharacterResponse
	CharacterSkillsResponse `json:"CharacterSkillsResponse"`
	Location                int64                       `json:"Location"`
	SkillQueue              []SkillQueue                `json:"SkillQueue"`
	QualifiedPlans          map[string]bool             `json:"QualifiedPlans"`
	PendingPlans            map[string]bool             `json:"PendingPlans"`
	PendingFinishDates      map[string]*time.Time       `json:"PendingFinishDates"`
	MissingSkills           map[string]map[string]int32 `json:"MissingSkills"`
}

// BaseCharacterResponse represents the user information returned by the EVE SSO
type BaseCharacterResponse struct {
	CharacterID   int64  `json:"CharacterID"`
	CharacterName string `json:"CharacterName"`
}

type CharacterResponse struct {
	AllianceID     int32     `json:"alliance_id,omitempty"`
	Birthday       time.Time `json:"birthday"`
	BloodlineID    int32     `json:"bloodline_id"`
	CorporationID  int32     `json:"corporation_id"`
	Description    string    `json:"description,omitempty"`
	FactionID      int32     `json:"faction_id,omitempty"`
	Gender         string    `json:"gender"`
	Name           string    `json:"name"`
	RaceID         int32     `json:"race_id"`
	SecurityStatus float64   `json:"security_status,omitempty"`
	Title          string    `json:"title,omitempty"`
}

type CharacterSkillsResponse struct {
	Skills        []SkillResponse `json:"skills"`
	TotalSP       int64           `json:"total_sp"`
	UnallocatedSP int32           `json:"unallocated_sp"`
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

type Identities struct {
	MainIdentity string                 `json:"main_identity"`
	Tokens       map[int64]oauth2.Token `json:"identities"`
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
