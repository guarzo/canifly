// persist/state.go
package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
)

// AppStateStore in-memory AppState protected by a mutex.
type AppStateStore struct {
	appState model.AppState
}

func (ds *DataStore) SetAppState(homeData model.AppState) {
	ds.appStateStoreLock(true)
	defer ds.appStateStoreUnlock(true)
	ds.appStateStore.appState = homeData
}

func (ds *DataStore) GetAppState() model.AppState {
	ds.appStateStoreLock(false)
	defer ds.appStateStoreUnlock(false)
	return ds.appStateStore.appState
}

func (ds *DataStore) SetAppStateLogin(isLoggedIn bool) error {
	appState := ds.GetAppState()
	appState.LoggedIn = isLoggedIn
	ds.SetAppState(appState)
	return ds.SaveAppStateSnapshot(appState)
}

func (ds *DataStore) ClearAppState() {
	ds.appStateStoreLock(true)
	defer ds.appStateStoreUnlock(true)
	ds.appStateStore.appState = model.AppState{}
}

// SaveAppStateSnapshot saves AppState to a JSON file.
func (ds *DataStore) SaveAppStateSnapshot(appState model.AppState) error {
	configDir, err := ds.GetWriteablePath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	snapshotPath := filepath.Join(configDir, "appstate_snapshot.json")

	data, err := json.MarshalIndent(appState, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal appState: %w", err)
	}

	if err := os.WriteFile(snapshotPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

// loadAppStateFromFile attempts to load AppState from disk.
func (ds *DataStore) loadAppStateFromFile() error {
	path, err := ds.getAppStateFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File does not exist, not an error, just no previous state
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read app state file: %w", err)
	}

	var appState model.AppState
	if err = json.Unmarshal(data, &appState); err != nil {
		return fmt.Errorf("failed to unmarshal app state: %w", err)
	}

	ds.appStateStoreLock(true)
	ds.appStateStore.appState = appState
	ds.appStateStoreUnlock(true)
	ds.logger.Debugf("Loaded persisted AppState from %s", path)

	return nil
}

func (ds *DataStore) getAppStateFilePath() (string, error) {
	configDir, err := ds.GetWriteablePath()
	if err != nil {
		return "", fmt.Errorf("failed to get config path: %w", err)
	}
	return filepath.Join(configDir, "appstate_snapshot.json"), nil
}

// Mutex methods for locking/unlocking the appState.
func (ds *DataStore) appStateStoreLock(write bool) {
	if write {
		ds.mut.Lock()
	} else {
		ds.mut.RLock()
	}
}

func (ds *DataStore) appStateStoreUnlock(write bool) {
	if write {
		ds.mut.Unlock()
	} else {
		ds.mut.RUnlock()
	}
}
