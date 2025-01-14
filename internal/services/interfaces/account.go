// services/interfaces/account.go
package interfaces

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

type AccountService interface {
	FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error
	UpdateAccountName(accountID int64, accountName string) error
	ToggleAccountStatus(accountID int64) error
	ToggleAccountVisibility(accountID int64) error
	RemoveAccountByName(accountName string) error
	RefreshAccountData(characterService CharacterService) (*model.AccountData, error)
	DeleteAllAccounts() error
	FetchAccounts() ([]model.Account, error)
	SaveAccounts(accounts []model.Account) error
	GetAccountNameByID(id string) (string, bool)
}
type AccountDataRepository interface {
	// FetchAccountData retrieves the entire account domain data (Accounts, UserAccount map, and Associations).
	FetchAccountData() (model.AccountData, error)

	// SaveAccountData saves the entire account domain data structure.
	SaveAccountData(data model.AccountData) error

	// DeleteAccountData removes the persisted account data file, if any.
	DeleteAccountData() error

	// FetchAccounts returns only the Accounts slice from the account data.
	// This is a convenience method that internally fetches AccountData and returns AccountData.Accounts.
	FetchAccounts() ([]model.Account, error)

	// SaveAccounts updates the Accounts slice in the account data, leaving UserAccount and Associations unchanged.
	SaveAccounts(accounts []model.Account) error

	// DeleteAccounts clears out the Accounts slice (but not necessarily UserAccount or Associations).
	DeleteAccounts() error
}

type AssociationService interface {
	UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error
	AssociateCharacter(userId, charId string) error
	UnassociateCharacter(userId, charId string) error
}

type CacheService interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadCache() error
	SaveCache() error
}

type CacheRepository interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadApiCache() error
	SaveApiCache() error
}

type CharacterService interface {
	ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
	DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
	UpdateCharacterFields(characterID int64, updates map[string]interface{}) error
	RemoveCharacter(characterID int64) error
}

// AuthClient defines methods related to authentication and token management.
type AuthClient interface {
	RefreshToken(refreshToken string) (*oauth2.Token, error)
	GetAuthURL(state string) string
	ExchangeCode(code string) (*oauth2.Token, error)
}

type EsiHttpClient interface {
	GetJSON(endpoint string, token *oauth2.Token, useCache bool, target interface{}) error
	GetJSONFromURL(url string, token *oauth2.Token, useCache bool, target interface{}) error
}

type DeletedCharactersRepository interface {
	FetchDeletedCharacters() ([]string, error)
	SaveDeletedCharacters([]string) error
}

type LoginService interface {
	ResolveAccountAndStatusByState(state string) (string, bool, bool)
	GenerateAndStoreInitialState(value string) (string, error)
	UpdateStateStatusAfterCallBack(state string) error
	ClearState(state string)
}

type LoginRepository interface {
	Set(state string, authStatus *model.AuthStatus)
	Get(state string) (*model.AuthStatus, bool)
	Delete(state string)
}

type SessionService interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
}

type SessionRepo interface {
	Get()
}
