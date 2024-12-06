// http/session.go
package http

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	AllAuthenticatedCharacters = "authenticated_characters"
	LoggedInUser               = "logged_in_user"
	SessionName                = "session"
	PreviousUserCount          = "previous_user_count"
	PreviousEtagUsed           = "previous_etag_used"
)

type SessionValues struct {
	LastRefreshTime   int64
	LoggedInUser      int64
	PreviousUserCount int
	PreviousEtagUsed  string
}

type SessionService struct {
	store *sessions.CookieStore
}

func GetSessionValues(session *sessions.Session) SessionValues {
	s := SessionValues{}

	if val, ok := session.Values[LoggedInUser].(int64); ok {
		s.LoggedInUser = val
	}

	if val, ok := session.Values[PreviousUserCount].(int); ok {
		s.PreviousUserCount = val
	}

	if val, ok := session.Values[PreviousEtagUsed].(string); ok {
		s.PreviousEtagUsed = val
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
