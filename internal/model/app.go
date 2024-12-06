package model

import (
	"encoding/gob"
	"time"

	"golang.org/x/oauth2"
)

type Account struct {
	Name       string
	Status     string // "Omega", "Alpha"
	Characters []CharacterIdentity
	ID         int64 // dynamically generated for now
}

type CharacterIdentity struct {
	Token     oauth2.Token // existing field
	Character Character    // existing field
	Role      string       // new field for categorization
	MCT       bool         // new field, useful for sorting
}

type UIData struct {
	Title      string
	LoggedIn   bool
	Accounts   []Account // Add Accounts here
	SkillPlans map[string]SkillPlanWithStatus
	ConfigData
}

type ConfigData struct {
	Roles []string `json:"Roles"`
}

func init() {
	gob.Register(CharacterIdentity{})
	gob.Register([]CharacterIdentity{})
	gob.Register([]Account{})
	gob.Register(Account{})
	gob.Register(Character{})
	gob.Register(UserInfoResponse{})
	gob.Register(CharacterSkillsResponse{})
	gob.Register(map[string]bool{})
	gob.Register(map[string]*time.Time{})
	gob.Register(map[string]map[string]int32{})
	gob.Register([]SkillQueue{})

}
