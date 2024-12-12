// model/config.go
package model

import (
	"encoding/gob"
	"time"
)

// AppState is the data passed to the UI
type AppState struct {
	LoggedIn    bool        `json:"LoggedIn"`
	AccountData AccountData `json:"AccountData"`
	ConfigData  ConfigData  `json:"ConfigData"`
	EveData     EveData     `json:"EveData"`
}

// DropDownSelections  are the dropdown selections on the sync page
type DropDownSelections map[string]UserSelection

type UserSelection struct {
	CharId string `json:"charId"`
	UserId string `json:"userId"`
}

// ConfigData are user settings and other app specific configuration
type ConfigData struct {
	Roles              []string `json:"Roles"`         // in app created roles for organizing data
	SettingsDir        string   `json:"SettingsDir"`   // directory where the settings are kept
	LastBackupDir      string   `json:"LastBackupDir"` // directory used for the previous backup
	DropDownSelections          // dropdown selections within the app
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
	gob.Register(AppState{})
	gob.Register(AccountData{})
	gob.Register(EveData{})
	gob.Register(EveProfile{})
	gob.Register(ConfigData{})
	gob.Register(DropDownSelections{})
	gob.Register([]EveProfile{})
	gob.Register([]Association{})
	gob.Register(Association{})

}
