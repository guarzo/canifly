// services/state/appstate_service.go
package config

import (
	"fmt"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.AppStateService = (*appStateService)(nil)

type appStateService struct {
	logger    interfaces.Logger
	stateRepo interfaces.AppStateRepository
}

func NewAppStateService(logger interfaces.Logger, ds interfaces.AppStateRepository) interfaces.AppStateService {
	return &appStateService{
		logger:    logger,
		stateRepo: ds,
	}
}

func (s *appStateService) GetAppState() model.AppState {
	return s.stateRepo.GetAppState()
}

func (s *appStateService) SetAppStateLogin(isLoggedIn bool) error {
	if err := s.stateRepo.SetAppStateLogin(isLoggedIn); err != nil {
		return fmt.Errorf("failed to set login state: %w", err)
	}
	return nil
}

func (s *appStateService) UpdateAndSaveAppState(data model.AppState) error {
	s.stateRepo.SetAppState(data)
	if err := s.stateRepo.SaveAppStateSnapshot(data); err != nil {
		return fmt.Errorf("failed to save app state snapshot: %w", err)
	}
	return nil
}

func (s *appStateService) ClearAppState() {
	s.stateRepo.ClearAppState()
}
