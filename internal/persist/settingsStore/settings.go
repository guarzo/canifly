// persist/datastore.go
package settingsStore

import (
	"sync"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.SettingsRepository = (*SettingsStore)(nil)

type SettingsStore struct {
	logger        interfaces.Logger
	appStateStore AppStateStore
	mut           sync.RWMutex
}

// NewConfigStore initializes a new SettingsStore and loads AppState from file.
func NewConfigStore(logger interfaces.Logger) *SettingsStore {
	s := &SettingsStore{
		logger:        logger,
		appStateStore: AppStateStore{},
	}

	if err := s.loadAppStateFromFile(); err != nil {
		s.logger.Warnf("unable to load app state from file %v", err)
	}

	return s
}
