package persist

import "sync"

type LoginStateStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewLoginStateStore() *LoginStateStore {
	return &LoginStateStore{
		store: make(map[string]string),
	}
}

func (l *LoginStateStore) Set(state, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.store[state] = value
}

func (l *LoginStateStore) Get(state string) (string, bool) {
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
