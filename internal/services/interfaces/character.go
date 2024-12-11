package interfaces

import "github.com/guarzo/canifly/internal/model"

type CharacterService interface {
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
	DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
	UpdateCharacterFields(characterID int64, updates map[string]interface{}) error
	RemoveCharacter(characterID int64) error
}

type DeletedCharactersRepository interface {
	FetchDeletedCharacters() ([]string, error)
	SaveDeletedCharacters([]string) error
}
