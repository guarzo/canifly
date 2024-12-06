package esi

import (
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

// ESIService defines the methods that our application needs from the ESI layer.
type ESIService interface {
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
}

type HTTPClient interface {
	DoRequest(method, endpoint string, body interface{}, target interface{}) error
}
