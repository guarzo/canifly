package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.AppStateRepository = (*AppStateStore)(nil)

type AppStateStore struct {
	logger   interfaces.Logger
	fs       persist.FileSystem
	basePath string

	mut      sync.RWMutex
	appState model.AppState
}

func NewAppStateStore(logger interfaces.Logger, fs persist.FileSystem, basePath string) *AppStateStore {
	store := &AppStateStore{
		logger:   logger,
		fs:       fs,
		basePath: basePath,
	}
	if err := store.loadAppStateFromFile(); err != nil {
		logger.Warnf("Unable to load app state from file: %v", err)
	}
	return store
}

func (s *AppStateStore) GetAppState() model.AppState {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.appState
}

func (s *AppStateStore) SetAppState(appState model.AppState) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.appState = appState
}

func (s *AppStateStore) SetAppStateLogin(isLoggedIn bool) error {
	appState := s.GetAppState()
	appState.LoggedIn = isLoggedIn
	s.SetAppState(appState)
	return s.SaveAppStateSnapshot(appState)
}

func (s *AppStateStore) ClearAppState() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.appState = model.AppState{}
}

func (s *AppStateStore) SaveAppStateSnapshot(appState model.AppState) error {
	snapshotPath := filepath.Join(s.basePath, "appstate_snapshot.json")
	s.logger.Debugf("app state saved at %s", snapshotPath)
	return persist.SaveJsonToFile(s.fs, snapshotPath, appState)
}

func (s *AppStateStore) loadAppStateFromFile() error {
	path := filepath.Join(s.basePath, "appstate_snapshot.json")
	if _, err := s.fs.Stat(path); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat appstate file: %w", err)
	}

	var appState model.AppState
	if err := persist.ReadJsonFromFile(s.fs, path, &appState); err != nil {
		return fmt.Errorf("failed to load AppState: %w", err)
	}

	s.mut.Lock()
	s.appState = appState
	s.mut.Unlock()
	s.logger.Debugf("Loaded persisted AppState from %s", path)
	return nil
}
