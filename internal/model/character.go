package model

import (
	"time"

	"golang.org/x/oauth2"
)

type AccountStatus string

const (
	Alpha AccountStatus = "Alpha"
	Omega AccountStatus = "Omega"
)

type Account struct {
	Name       string
	Status     AccountStatus
	Characters []CharacterIdentity
	ID         int64 // dynamically generated for now
}

type CharacterIdentity struct {
	Token     oauth2.Token
	Character Character
	Role      string
	MCT       bool
	Training  string
}

type Character struct {
	UserInfoResponse
	CharacterSkillsResponse `json:"CharacterSkillsResponse"`
	Location                int64  `json:"Location"`
	LocationName            string `json:"LocationName"`

	SkillQueue         []SkillQueue                `json:"SkillQueue"`
	QualifiedPlans     map[string]bool             `json:"QualifiedPlans"`
	PendingPlans       map[string]bool             `json:"PendingPlans"`
	PendingFinishDates map[string]*time.Time       `json:"PendingFinishDates"`
	MissingSkills      map[string]map[string]int32 `json:"MissingSkills"`
}

// UserInfoResponse represents the user information returned by the EVE SSO
type UserInfoResponse struct {
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
