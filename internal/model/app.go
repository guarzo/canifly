package model

import (
	"encoding/gob"
	"time"
)

type AppState struct {
	Title      string
	LoggedIn   bool
	Accounts   []Account // Add Accounts here
	SkillPlans map[string]SkillPlanWithStatus
	ConfigData
	SubDirs []SubDirData
	UserSelections
}

type UserSelection struct {
	CharId string `json:"charId"`
	UserId string `json:"userId"`
}

type UserSelections map[string]UserSelection

type SubDirData struct {
	SubDir             string     `json:"subDir"`
	AvailableCharFiles []CharFile `json:"availableCharFiles"`
	AvailableUserFiles []UserFile `json:"availableUserFiles"`
}

type ConfigData struct {
	Roles        []string      `json:"Roles"`
	SettingsDir  string        `json:"SettingsDir"`
	SettingsData []SubDirData  `json:"settingsData"`
	Associations []Association `json:"associations"`
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

type Association struct {
	UserId      string `json:"userId"`
	CharId      string `json:"charId"`
	CharName    string `json:"charName"`
	AccountName string `json:"accountName,omitempty"`
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
