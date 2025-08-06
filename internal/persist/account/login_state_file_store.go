package account

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.LoginRepository = (*LoginStateFileStore)(nil)

type LoginStateFileStore struct {
	mu       sync.RWMutex
	store    map[string]*model.AuthStatus
	filePath string
	logger   interfaces.Logger
	fs       persist.FileSystem
}

type persistedAuthStatus struct {
	*model.AuthStatus
	Timestamp time.Time `json:"timestamp"`
}

func NewLoginStateFileStore(basePath string, logger interfaces.Logger) (*LoginStateFileStore, error) {
	filePath := filepath.Join(basePath, "login_states.json")

	store := &LoginStateFileStore{
		store:    make(map[string]*model.AuthStatus),
		filePath: filePath,
		logger:   logger,
		fs:       persist.NewOSFileSystem(),
	}

	// Load existing states from file
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		logger.Warnf("Failed to load login states: %v", err)
	}

	// Clean up old states periodically (older than 1 hour)
	go store.cleanupOldStates()

	return store, nil
}

func (l *LoginStateFileStore) Set(state string, authStatus *model.AuthStatus) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Infof("Setting login state: %s -> account=%s, callback=%v", state, authStatus.AccountName, authStatus.CallBackComplete)
	l.store[state] = authStatus
	if err := l.save(); err != nil {
		l.logger.Errorf("Failed to save login states: %v", err)
	} else {
		l.logger.Debugf("Successfully saved login state to file: %s", l.filePath)
	}
}

func (l *LoginStateFileStore) Get(state string) (*model.AuthStatus, bool) {
	// Need write lock since we're reloading from file
	l.mu.Lock()
	defer l.mu.Unlock()

	// Reload from file to get latest state
	if err := l.loadWithoutLock(); err != nil && !os.IsNotExist(err) {
		l.logger.Warnf("Failed to reload login states: %v", err)
	}

	val, ok := l.store[state]
	if ok {
		l.logger.Infof("Found login state: %s -> account=%s, callback=%v", state, val.AccountName, val.CallBackComplete)
	} else {
		l.logger.Warnf("Login state not found: %s. Available states: %d", state, len(l.store))
		// Log available states for debugging
		for s := range l.store {
			l.logger.Debugf("  Available state: %s", s)
		}
	}
	return val, ok
}

func (l *LoginStateFileStore) Delete(state string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.store, state)
	if err := l.save(); err != nil {
		l.logger.Errorf("Failed to save login states after delete: %v", err)
	}
}

func (l *LoginStateFileStore) load() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.loadWithoutLock()
}

func (l *LoginStateFileStore) loadWithoutLock() error {
	data, err := l.fs.ReadFile(l.filePath)
	if err != nil {
		return err
	}

	var persisted map[string]*persistedAuthStatus
	if err := json.Unmarshal(data, &persisted); err != nil {
		return err
	}

	// Convert back to regular AuthStatus and filter old entries
	now := time.Now()
	for state, p := range persisted {
		// Skip states older than 1 hour
		if now.Sub(p.Timestamp) > time.Hour {
			continue
		}
		l.store[state] = p.AuthStatus
	}

	return nil
}

func (l *LoginStateFileStore) save() error {
	// Convert to persisted format with timestamps
	persisted := make(map[string]*persistedAuthStatus)
	now := time.Now()
	for state, auth := range l.store {
		persisted[state] = &persistedAuthStatus{
			AuthStatus: auth,
			Timestamp:  now,
		}
	}

	data, err := json.MarshalIndent(persisted, "", "  ")
	if err != nil {
		return err
	}

	return persist.AtomicWriteFile(l.fs, l.filePath, data, 0600)
}

func (l *LoginStateFileStore) cleanupOldStates() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()

		// Reload and clean
		if err := l.loadWithoutLock(); err != nil && !os.IsNotExist(err) {
			l.logger.Warnf("Failed to reload during cleanup: %v", err)
		}

		// Save cleaned state
		if err := l.save(); err != nil {
			l.logger.Errorf("Failed to save during cleanup: %v", err)
		}

		l.mu.Unlock()
	}
}
