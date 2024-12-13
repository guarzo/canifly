package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/testutil"
)

const (
	SessionName = "testSession"
	LoggedIn    = "logged_in"
)

// MockSessionService simulates the behavior of SessionService.
// Now it returns a real *sessions.Session instead of a mock session.
type MockSessionService struct {
	Store    *sessions.CookieStore
	Err      error
	LoggedIn bool
}

func (m *MockSessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.Store == nil {
		m.Store = sessions.NewCookieStore([]byte("secret"))
	}

	// Get the session once
	session, _ := m.Store.Get(r, name)

	// If we want the user to be logged in, set the value here
	if m.LoggedIn {
		session.Values[LoggedIn] = true
	}

	// Return the modified session
	return session, nil
}

// createTestRouter creates a router with the AuthMiddleware applied and test handlers.
func createTestRouter(sessionService interfaces.SessionService, logger interfaces.Logger) *mux.Router {
	r := mux.NewRouter()
	r.Use(flyHttp.AuthMiddleware(sessionService, logger))

	// Public routes
	r.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "public ok"})
	})
	r.HandleFunc("/landing", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "landing ok"})
	})

	// Private routes
	r.HandleFunc("/private", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "private ok"})
	})

	return r
}

func TestAuthMiddleware_PublicRoute(t *testing.T) {
	sessionService := &MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	router := createTestRouter(sessionService, logger)

	req, _ := http.NewRequest("GET", "/static", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "public ok", resp["status"])
}

func TestAuthMiddleware_PublicRoute_Landing(t *testing.T) {
	sessionService := &MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	router := createTestRouter(sessionService, logger)

	req, _ := http.NewRequest("GET", "/landing", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "landing ok", resp["status"])
}

func TestAuthMiddleware_PrivateRoute_NotLoggedIn(t *testing.T) {
	sessionService := &MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	router := createTestRouter(sessionService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "user is not logged in", resp["error"])
}

func TestAuthMiddleware_PrivateRoute_SessionError(t *testing.T) {
	sessionService := &MockSessionService{
		Err: errors.New("session retrieval error"),
	}
	logger := &testutil.MockLogger{}

	router := createTestRouter(sessionService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "failed to retrieve session", resp["error"])
}

func TestAuthMiddleware_PrivateRoute_LoggedIn(t *testing.T) {
	sessionService := &MockSessionService{
		Store:    sessions.NewCookieStore([]byte("secret")),
		LoggedIn: true,
	}
	logger := &testutil.MockLogger{}

	router := createTestRouter(sessionService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Should be authenticated")
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "private ok", resp["status"])
}
