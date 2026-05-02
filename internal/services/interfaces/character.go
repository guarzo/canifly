package interfaces

import "github.com/guarzo/canifly/internal/model"

// CharacterService handles character management operations
type CharacterService interface {
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
	DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
	UpdateCharacter(characterID int64, update model.CharacterUpdate) error
	RemoveCharacter(characterID int64) error
	RefreshCharacterData(characterID int64) (bool, error)
}
