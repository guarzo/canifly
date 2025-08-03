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
	store := sessions.NewCookieStore([]byte(secret))

	// Configure cookie settings for CORS
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to false for localhost development
		SameSite: http.SameSiteLaxMode,
	}

	return &sessionService{
		store: store,
	}
}

func (s *sessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.store.Get(r, name)
}
