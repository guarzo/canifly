package model

import (
	"time"

	"golang.org/x/oauth2"
)

type HomeData struct {
	Title               string
	LoggedIn            bool
	Identities          map[int64]CharacterData
	TabulatorIdentities []map[string]interface{}
	TabulatorSkillPlans map[string]SkillPlan
	MainIdentity        int64
	MatchingSkillPlans  map[string]SkillPlan
	MatchingCharacters  map[int64]CharacterData
}

// CharacterData structure contains nested Character struct
type CharacterData struct {
	Token     oauth2.Token
	Character Character // Keeping Character as an independent struct
}

// Character represents the user information
type Character struct {
	User
	CharacterSkillsResponse `json:"CharacterSkillsResponse"`
	SkillQueue              []SkillQueue                `json:"SkillQueue"`
	QualifiedPlans          map[string]bool             `json:"QualifiedPlans"`
	PendingPlans            map[string]bool             `json:"PendingPlans"`
	PendingFinishDates      map[string]*time.Time       `json:"PendingFinishDates"`
	MissingSkills           map[string]map[string]int32 `json:"MissingSkills"`
}

// User represents the user information returned by the EVE SSO
type User struct {
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
