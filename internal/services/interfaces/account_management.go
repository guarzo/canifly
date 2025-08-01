package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
	"golang.org/x/oauth2"
)

// AccountManagementService consolidates AccountService and AssociationService interfaces
type AccountManagementService interface {
	// Account Management (from AccountService)
	FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error
	UpdateAccountName(accountID int64, accountName string) error
	ToggleAccountStatus(accountID int64) error
	ToggleAccountVisibility(accountID int64) error
	RemoveAccountByName(accountName string) error
	RefreshAccountData(eveDataService EVEDataService) (*model.AccountData, error)
	DeleteAllAccounts() error
	FetchAccounts() ([]model.Account, error)
	SaveAccounts(accounts []model.Account) error
	GetAccountNameByID(id string) (string, bool)

	// Association Management (from AssociationService)
	UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error
	AssociateCharacter(userId, charId string) error
	UnassociateCharacter(userId, charId string) error
}