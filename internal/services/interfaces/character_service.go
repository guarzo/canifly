package interfaces

import "github.com/guarzo/canifly/internal/model"

type CharacterService interface {
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
}
