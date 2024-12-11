// persist/state.go
package settingsStore

import (
	"fmt"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

// AppStateStore in-memory AppState protected by a mutex.
type AppStateStore struct {
	appState model.AppState
}

func (s *SettingsStore) SetAppState(homeData model.AppState) {
	s.appStateStoreLock(true)
	defer s.appStateStoreUnlock(true)
	s.appStateStore.appState = homeData
}

func (s *SettingsStore) GetAppState() model.AppState {
	s.appStateStoreLock(false)
	defer s.appStateStoreUnlock(false)
	return s.appStateStore.appState
}

func (s *SettingsStore) SetAppStateLogin(isLoggedIn bool) error {
	appState := s.GetAppState()
	appState.LoggedIn = isLoggedIn
	s.SetAppState(appState)
	return s.SaveAppStateSnapshot(appState)
}

func (s *SettingsStore) ClearAppState() {
	s.appStateStoreLock(true)
	defer s.appStateStoreUnlock(true)
	s.appStateStore.appState = model.AppState{}
}

// SaveAppStateSnapshot saves AppState to a JSON file.
func (s *SettingsStore) SaveAppStateSnapshot(appState model.AppState) error {
	snapshotPath, err := s.getAppStateFilePath()
	if err != nil {
		return err
	}

	s.logger.Infof("app state saved at %s", snapshotPath)

	return persist.OldSaveJson(snapshotPath, appState)
}

// loadAppStateFromFile attempts to load AppState from disk.
func (s *SettingsStore) loadAppStateFromFile() error {
	path, err := s.getAppStateFilePath()
	if err != nil {
		return err
	}
	if _, err := persist.FileExists(path); err != nil {
		return nil
	}

	var appState model.AppState
	if err = persist.OldReadJson(path, &appState); err != nil {
		return fmt.Errorf("failed to load AppState: %w", err)
	}

	s.appStateStoreLock(true)
	s.appStateStore.appState = appState
	s.appStateStoreUnlock(true)
	s.logger.Debugf("Loaded persisted AppState from %s", path)

	return nil
}

func (s *SettingsStore) getAppStateFilePath() (string, error) {
	return getConfigFileName("appstate_snapshot.json")
}

// Mutex methods for locking/unlocking the appState.
func (s *SettingsStore) appStateStoreLock(write bool) {
	if write {
		s.mut.Lock()
	} else {
		s.mut.RLock()
	}
}

func (s *SettingsStore) appStateStoreUnlock(write bool) {
	if write {
		s.mut.Unlock()
	} else {
		s.mut.RUnlock()
	}
}
