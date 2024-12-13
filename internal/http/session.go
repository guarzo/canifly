// Package http - http/session.go
package http

import (
	"github.com/guarzo/canifly/internal/services/interfaces"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "session"
	LoggedIn    = "logged_in"
)

type sessionService struct {
	store *sessions.CookieStore
}

func NewSessionService(secret string) interfaces.SessionService {
	return &sessionService{
		store: sessions.NewCookieStore([]byte(secret)),
	}
}

func (s *sessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.store.Get(r, name)
}
