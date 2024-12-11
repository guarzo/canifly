package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
	"golang.org/x/oauth2"
)

type ESIService interface {
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacter(id string) (*model.CharacterResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
	ResolveCharacterNames(charIds []string) (map[string]string, error)
	SaveEsiCache() error
}

// AuthClient defines methods related to authentication and token management.
type AuthClient interface {
	// RefreshToken takes a refresh token and returns a new, updated oauth2.Token
	RefreshToken(refreshToken string) (*oauth2.Token, error)
}

type HTTPClient interface {
	// DoRequest executes an HTTP request with the given method, endpoint, and optional body,
	// then unmarshals the result into target.
	DoRequest(method, endpoint string, body interface{}, target interface{}) error
}

type DeletedCharactersRepository interface {
	FetchDeletedCharacters() ([]string, error)
	SaveDeletedCharacters([]string) error
}
