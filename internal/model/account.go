// model/account.go
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
	ID         int64 // userFile ID for this account, defaults to 0 until assigned
	Visible    bool  // toggle visibility
}

type AccountData struct {
	Accounts     []Account
	Associations []Association // in app assigned connections between user and char files (effectively connecting characters to accounts)
}

// Association are the user connections between char and user files (userID does correspond to accountId)
type Association struct {
	UserId   string `json:"userId"`
	CharId   string `json:"charId"`
	CharName string `json:"charName"`
}

type CharacterIdentity struct {
	Token           oauth2.Token
	Character       Character
	CorporationName string
	AllianceName    string
	Role            string
	MCT             bool
	Training        string
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

type AuthStatus struct {
	AccountName      string `json:"accountName"`
	CallBackComplete bool   `json:"callBackComplete"`
}

// Alliance contains detailed information about an EVE Online alliance
type Alliance struct {
	CreatorCorporationID  int       `json:"creator_corporation_id"`
	CreatorID             int       `json:"creator_id"`
	DateFounded           time.Time `json:"date_founded"`
	ExecutorCorporationID int       `json:"executor_corporation_id"`
	Name                  string    `json:"name"`
	Ticker                string    `json:"ticker"`
}

// Corporation represents detailed information about an EVE Online corporation
type Corporation struct {
	AllianceID    int       `json:"alliance_id"`
	CeoID         int       `json:"ceo_id"`
	CreatorID     int       `json:"creator_id"`
	DateFounded   time.Time `json:"date_founded"`
	Description   string    `json:"description"`
	HomeStationID int       `json:"home_station_id"`
	MemberCount   int       `json:"member_count"`
	Name          string    `json:"name"`
	Shares        int       `json:"shares"`
	TaxRate       float64   `json:"tax_rate"`
	Ticker        string    `json:"ticker"`
	URL           string    `json:"url"`
}
