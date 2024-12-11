package interfaces

import "github.com/guarzo/canifly/internal/model"

type AssociationService interface {
	UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error
	AssociateCharacter(userId, charId string) error
	UnassociateCharacter(userId, charId string) error
}

type AssocRepository interface {
	ConfigRepository
	AccountRepository
}
