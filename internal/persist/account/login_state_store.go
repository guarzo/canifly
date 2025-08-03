package account

import (
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.LoginRepository = (*LoginStateStore)(nil)

type LoginStateStore struct {
	mu    sync.Mutex
	store map[string]*model.AuthStatus
}

func NewLoginStateStore() *LoginStateStore {
	return &LoginStateStore{
		store: make(map[string]*model.AuthStatus),
	}
}

func (l *LoginStateStore) Set(state string, authStatus *model.AuthStatus) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.store[state] = authStatus
}

func (l *LoginStateStore) Get(state string) (*model.AuthStatus, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	val, ok := l.store[state]
	return val, ok
}

func (l *LoginStateStore) Delete(state string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.store, state)
}
