// persist/snapshot.go
package persist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
)

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

	if err := ioutil.WriteFile(snapshotPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

// loadAppStateFromFile attempts to load AppState from a JSON file on disk.
func (ds *DataStore) loadAppStateFromFile() error {
	path, _ := ds.getAppStateFilePath()
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

	// If successfully unmarshaled, set it and compute ETag
	etag, err := GenerateETag(appState)
	if err != nil {
		return fmt.Errorf("failed to generate ETag: %w", err)
	}

	ds.appStateStore.Lock()
	defer ds.appStateStore.Unlock()
	ds.appStateStore.appState = appState
	ds.appStateStore.ETag = etag
	ds.logger.Infof("Loaded persisted AppState from %s", path)

	return nil
}

func (ds *DataStore) getAppStateFilePath() (string, error) {
	configDir, err := ds.GetWriteablePath()
	if err != nil {
		return "", fmt.Errorf("failed to get config path: %w", err)
	}

	return filepath.Join(configDir, "appstate_snapshot.json"), nil

}
