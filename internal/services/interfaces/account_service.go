package interfaces

import (
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

type AccountService interface {
	FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error
	AssociateCharacter(userId, charId string) error
	UnassociateCharacter(userId, charId string) error
	UpdateAccountName(accountID int64, accountName string) error
	ToggleAccountStatus(accountID int64) error
	RemoveAccountByName(accountName string) error
	RefreshAccounts(characterService CharacterService) ([]model.Account, error)
}
