package interfaces

import "github.com/guarzo/canifly/internal/model"

// AccountDataProvider provides access to account data
type AccountDataProvider interface {
	FetchAccounts() ([]model.Account, error)
	SaveAccounts(accounts []model.Account) error
	GetAccountNameByID(id string) (string, bool)
}