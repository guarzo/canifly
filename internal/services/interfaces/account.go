package interfaces

import (
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

type AccountService interface {
	FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error
	UpdateAccountName(accountID int64, accountName string) error
	ToggleAccountStatus(accountID int64) error
	RemoveAccountByName(accountName string) error
	RefreshAccounts(characterService CharacterService) ([]model.Account, error)
	DeleteAllAccounts() error
	FetchAccounts() ([]model.Account, error)
	SaveAccounts(accounts []model.Account) error
}

type AccountRepository interface {
	FetchAccounts() ([]model.Account, error)
	SaveAccounts([]model.Account) error
	DeleteAccounts() error
}
