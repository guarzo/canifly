package account

import (
	"errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.LoginService = (*loginService)(nil)

type loginService struct {
	logger    interfaces.Logger
	loginRepo interfaces.LoginRepository
}

func NewLoginService(logger interfaces.Logger, loginRepo interfaces.LoginRepository) interfaces.LoginService {
	return &loginService{
		logger:    logger,
		loginRepo: loginRepo,
	}
}

func (l *loginService) GenerateAndStoreInitialState(value string) (string, error) {
	state, err := persist.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	l.loginRepo.Set(state, &model.AuthStatus{
		AccountName:      value,
		CallBackComplete: false,
	})
	return state, nil
}

func (l *loginService) ResolveAccountAndStatusByState(state string) (string, bool, bool) {
	authStatus, ok := l.loginRepo.Get(state)
	if !ok {
		return "", false, false
	}
	return authStatus.AccountName, authStatus.CallBackComplete, true
}

func (l *loginService) UpdateStateStatusAfterCallBack(state string) error {
	authStatus, ok := l.loginRepo.Get(state)
	if !ok {
		return errors.New("unable to retrieve authStatus for provided state")
	}
	authStatus.CallBackComplete = true
	l.loginRepo.Set(state, authStatus)
	return nil
}

func (l *loginService) ClearState(state string) {
	l.loginRepo.Delete(state)
}
