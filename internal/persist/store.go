// persist/store.go
package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/guarzo/canifly/internal/model"
)

// DataStore now holds a global AppState in memory along with its ETag
// protected by a mutex.
type AppStateStore struct {
	sync.RWMutex
	appState model.AppState
	ETag     string
}

// GenerateETag computes a hash of the AppState
func GenerateETag(homeData model.AppState) (string, error) {
	data, err := json.Marshal(homeData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// SetAppState updates the in-memory AppState and computes a new ETag.
func (ds *DataStore) SetAppState(homeData model.AppState) (string, error) {
	ds.appStateStore.Lock()
	defer ds.appStateStore.Unlock()
	ds.appStateStore.appState = homeData

	etag, err := GenerateETag(homeData)
	if err != nil {
		return "", err
	}
	ds.appStateStore.ETag = etag
	return etag, nil
}

// GetAppState returns the currently stored AppState and its ETag.
// The boolean indicates if a state has been set yet (if empty, false).
func (ds *DataStore) GetAppState() (model.AppState, string, bool) {
	ds.appStateStore.RLock()
	defer ds.appStateStore.RUnlock()
	// Check if appState is "empty" by some logic, for example if no accounts and title empty
	// or if we just rely on the fact that ETag is empty if nothing stored yet.
	if ds.appStateStore.ETag == "" {
		return model.AppState{}, "", false
	}
	return ds.appStateStore.appState, ds.appStateStore.ETag, true
}

// ClearAppState removes the stored AppState
func (ds *DataStore) ClearAppState() {
	ds.appStateStore.Lock()
	defer ds.appStateStore.Unlock()
	ds.appStateStore.appState = model.AppState{}
	ds.appStateStore.ETag = ""
}
