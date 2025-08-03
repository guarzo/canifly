package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
	"golang.org/x/oauth2"
)

// ESIAPIService handles all EVE ESI API operations
type ESIAPIService interface {
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacter(id string) (*model.CharacterResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
	ResolveCharacterNames(charIds []string) (map[string]string, error)
	GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error)
	GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error)
}
