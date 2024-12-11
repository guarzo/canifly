package interfaces

import "github.com/guarzo/canifly/internal/model"

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
