// persist/state.go
package persist

import (
	"fmt"
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
	snapshotPath, err := ds.getAppStateFilePath()
	if err != nil {
		return err
	}

	ds.logger.Infof("app state saved at %s", snapshotPath)

	return saveJSONToFile(snapshotPath, appState)
}

// loadAppStateFromFile attempts to load AppState from disk.
func (ds *DataStore) loadAppStateFromFile() error {
	path, err := ds.getAppStateFilePath()
	if err != nil {
		return err
	}
	if _, err := fileExists(path); err != nil {
		return nil
	}

	var appState model.AppState
	if err = readJSONFromFile(path, &appState); err != nil {
		return fmt.Errorf("failed to load AppState: %w", err)
	}

	ds.appStateStoreLock(true)
	ds.appStateStore.appState = appState
	ds.appStateStoreUnlock(true)
	ds.logger.Debugf("Loaded persisted AppState from %s", path)

	return nil
}

func (ds *DataStore) getAppStateFilePath() (string, error) {
	return getConfigFileName("appstate_snapshot.json")
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
