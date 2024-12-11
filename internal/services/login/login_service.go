package login

import (
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/utils"
)

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

func (l *loginService) GenerateAnStoreState(value string) (string, error) {
	state, err := utils.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	l.loginRepo.Set(state, value)
	return state, nil
}

func (l *loginService) ResolveAccountByState(state string) (string, bool) {
	return l.loginRepo.Get(state)
}

func (l *loginService) ClearState(state string) {
	l.loginRepo.Delete(state)
}
