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

// MockLoginService is a mock implementation of LoginService for testing
type MockLoginService struct {
	states map[string]struct {
		accountName      string
		callbackComplete bool
	}
}

func NewMockLoginService() *MockLoginService {
	return &MockLoginService{
		states: make(map[string]struct {
			accountName      string
			callbackComplete bool
		}),
	}
}

func (m *MockLoginService) GenerateAndStoreInitialState(value string) (string, error) {
	return "test-state", nil
}

func (m *MockLoginService) ResolveAccountAndStatusByState(state string) (string, bool, bool) {
	if s, ok := m.states[state]; ok {
		return s.accountName, s.callbackComplete, true
	}
	return "", false, false
}

func (m *MockLoginService) UpdateStateStatusAfterCallBack(state string) error {
	if s, ok := m.states[state]; ok {
		s.callbackComplete = true
		m.states[state] = s
		return nil
	}
	return errors.New("state not found")
}

func (m *MockLoginService) ClearState(state string) {
	delete(m.states, state)
}

// stubPersistentValidator is a minimal PersistentSessionValidator for tests.
// validFn lets each test decide which tokens are considered valid; called is
// flipped on every invocation so tests can assert whether the persistent path
// was consulted.
type stubPersistentValidator struct {
	validFn func(string) bool
	called  bool
}

func (s *stubPersistentValidator) ValidateSession(sessionID string) bool {
	s.called = true
	if s.validFn == nil {
		return false
	}
	return s.validFn(sessionID)
}

// createTestRouter creates a router with the AuthMiddleware applied and test handlers.
func createTestRouter(sessionService interfaces.SessionService, loginService interfaces.LoginService, logger interfaces.Logger) *mux.Router {
	return createTestRouterWithValidator(sessionService, loginService, nil, logger)
}

func createTestRouterWithValidator(sessionService interfaces.SessionService, loginService interfaces.LoginService, validator flyHttp.PersistentSessionValidator, logger interfaces.Logger) *mux.Router {
	r := mux.NewRouter()
	r.Use(flyHttp.AuthMiddleware(sessionService, loginService, validator, logger))

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
	sessionService := &testutil.MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	loginService := NewMockLoginService()
	router := createTestRouter(sessionService, loginService, logger)

	req, _ := http.NewRequest("GET", "/static", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "public ok", resp["status"])
}

func TestAuthMiddleware_PublicRoute_Landing(t *testing.T) {
	sessionService := &testutil.MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	loginService := NewMockLoginService()
	router := createTestRouter(sessionService, loginService, logger)

	req, _ := http.NewRequest("GET", "/landing", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "landing ok", resp["status"])
}

func TestAuthMiddleware_PrivateRoute_NotLoggedIn(t *testing.T) {
	sessionService := &testutil.MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}

	loginService := NewMockLoginService()
	router := createTestRouter(sessionService, loginService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "user is not logged in", resp["error"])
}

func TestAuthMiddleware_PrivateRoute_SessionError(t *testing.T) {
	sessionService := &testutil.MockSessionService{
		Err: errors.New("session retrieval error"),
	}
	logger := &testutil.MockLogger{}

	loginService := NewMockLoginService()
	router := createTestRouter(sessionService, loginService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "failed to retrieve session", resp["error"])
}

func TestAuthMiddleware_PrivateRoute_LoggedIn(t *testing.T) {
	sessionService := &testutil.MockSessionService{
		Store:    sessions.NewCookieStore([]byte("secret")),
		LoggedIn: true,
	}
	logger := &testutil.MockLogger{}

	loginService := NewMockLoginService()
	router := createTestRouter(sessionService, loginService, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Should be authenticated")
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "private ok", resp["status"])
}

func TestAuthMiddleware_BearerToken_PersistentSessionValid(t *testing.T) {
	sessionService := &testutil.MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}
	loginService := NewMockLoginService()

	const token = "persistent-session-token"
	validator := &stubPersistentValidator{
		validFn: func(id string) bool { return id == token },
	}

	router := createTestRouterWithValidator(sessionService, loginService, validator, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.True(t, validator.called, "persistent validator should be consulted first")
	assert.Equal(t, http.StatusOK, rr.Code, "valid persistent token should authenticate")
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "private ok", resp["status"])
}

func TestAuthMiddleware_BearerToken_FallsBackToOAuthState(t *testing.T) {
	sessionService := &testutil.MockSessionService{Store: sessions.NewCookieStore([]byte("secret"))}
	logger := &testutil.MockLogger{}
	loginService := NewMockLoginService()

	// Seed a completed OAuth state to act as the legacy bearer token.
	const token = "legacy-oauth-state"
	loginService.states[token] = struct {
		accountName      string
		callbackComplete bool
	}{accountName: "acct-1", callbackComplete: true}

	validator := &stubPersistentValidator{
		validFn: func(string) bool { return false },
	}

	router := createTestRouterWithValidator(sessionService, loginService, validator, logger)

	req, _ := http.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.True(t, validator.called, "persistent validator should be consulted before fallback")
	assert.Equal(t, http.StatusOK, rr.Code, "should authenticate via OAuth state fallback")
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "private ok", resp["status"])

	// Sanity check: the fallback path is what resolved the token.
	acct, complete, ok := loginService.ResolveAccountAndStatusByState(token)
	assert.True(t, ok)
	assert.True(t, complete)
	assert.Equal(t, "acct-1", acct)
}
