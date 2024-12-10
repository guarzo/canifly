// services/state/state_service.go
package state

import (
	"fmt"
	"github.com/guarzo/canifly/internal/model"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type stateService struct {
	logger    interfaces.Logger
	dataStore interfaces.DataStore
}

func NewStateService(logger interfaces.Logger, ds interfaces.DataStore) interfaces.StateService {
	return &stateService{logger: logger, dataStore: ds}
}

func (s *stateService) SetAppStateLogin(isLoggedIn bool) error {
	if err := s.dataStore.SetAppStateLogin(isLoggedIn); err != nil {
		return fmt.Errorf("failed to set login state: %w", err)
	}
	return nil
}

func (s *stateService) DeleteAllAccounts() error {
	if err := s.dataStore.DeleteAccounts(); err != nil {
		return fmt.Errorf("failed to delete accounts: %w", err)
	}
	return nil
}

func (s *stateService) UpdateAndSaveAppState(data model.AppState) error {
	s.dataStore.SetAppState(data)
	if err := s.dataStore.SaveAppStateSnapshot(data); err != nil {
		return fmt.Errorf("failed to save app state snapshot: %w", err)
	}
	return nil
}
