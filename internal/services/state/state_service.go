// services/state/state_service.go
package state

import (
	"fmt"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type stateService struct {
	logger    interfaces.Logger
	stateRepo interfaces.StateRepository
}

func NewStateService(logger interfaces.Logger, ds interfaces.StateRepository) interfaces.StateService {
	return &stateService{
		logger:    logger,
		stateRepo: ds,
	}
}

func (s *stateService) GetAppState() model.AppState {
	return s.stateRepo.GetAppState()
}

func (s *stateService) SetAppStateLogin(isLoggedIn bool) error {
	if err := s.stateRepo.SetAppStateLogin(isLoggedIn); err != nil {
		return fmt.Errorf("failed to set login state: %w", err)
	}
	return nil
}

func (s *stateService) UpdateAndSaveAppState(data model.AppState) error {
	s.stateRepo.SetAppState(data)
	if err := s.stateRepo.SaveAppStateSnapshot(data); err != nil {
		return fmt.Errorf("failed to save app state snapshot: %w", err)
	}
	return nil
}
