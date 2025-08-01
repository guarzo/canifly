package interfaces

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

// Logger interface
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithError(err error) Logger
	WithFields(fields map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
}


// CacheService interface
type CacheService interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadCache() error
	SaveCache() error
}

// AuthClient interface
type AuthClient interface {
	RefreshToken(refreshToken string) (*oauth2.Token, error)
	GetAuthURL(state string) string
	ExchangeCode(code string) (*oauth2.Token, error)
}

// EsiHttpClient interface
type EsiHttpClient interface {
	GetJSON(endpoint string, token *oauth2.Token, useCache bool, target interface{}) error
	GetJSONFromURL(url string, token *oauth2.Token, useCache bool, target interface{}) error
}

// LoginService interface
type LoginService interface {
	ResolveAccountAndStatusByState(state string) (string, bool, bool)
	GenerateAndStoreInitialState(value string) (string, error)
	UpdateStateStatusAfterCallBack(state string) error
	ClearState(state string)
}

// SessionService interface
type SessionService interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
}