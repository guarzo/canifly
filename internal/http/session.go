// Package http - http/session.go
package http

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "session"
	LoggedIn    = "logged_in"
)

type SessionValues struct {
	LoggedIn bool
}

type SessionService struct {
	store *sessions.CookieStore
}

func GetSessionValues(session *sessions.Session) SessionValues {
	s := SessionValues{}

	if val, ok := session.Values[LoggedIn].(bool); ok {
		s.LoggedIn = val
	}

	return s
}

func NewSessionService(secret string) *SessionService {
	return &SessionService{
		store: sessions.NewCookieStore([]byte(secret)),
	}
}

func (s *SessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.store.Get(r, name)
}
