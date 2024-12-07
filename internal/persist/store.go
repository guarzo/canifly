// persist/store.go
package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/guarzo/canifly/internal/model"
)

var Store *UIDataStore

// UIDataStore stores AppState in memory
type UIDataStore struct {
	sync.RWMutex
	store map[int64]model.AppState
	ETag  string
}

// NewUIDataStore creates a new UIDataStore
func NewUIDataStore() *UIDataStore {
	return &UIDataStore{
		store: make(map[int64]model.AppState),
		ETag:  "",
	}
}

func (s *UIDataStore) Set(id int64, homeData model.AppState) (string, error) {
	s.Lock()
	defer s.Unlock()
	s.store[id] = homeData

	etag, err := GenerateETag(homeData)
	if err != nil {
		return "", err
	}
	s.ETag = etag
	return etag, nil
}

func (s *UIDataStore) Get(id int64) (model.AppState, string, bool) {
	s.RLock()
	defer s.RUnlock()
	homeData, ok := s.store[id]
	return homeData, s.ETag, ok
}

// Delete removes an identity from the persist
func (s *UIDataStore) Delete(id int64) {
	s.Lock()
	defer s.Unlock()
	delete(s.store, id)
}

func GenerateETag(homeData model.AppState) (string, error) {
	data, err := json.Marshal(homeData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func init() {
	Store = NewUIDataStore()
}
