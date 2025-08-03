package http_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/testutil"
)

func TestAPIClient_GetJSON_SuccessNoCache(t *testing.T) {
	// Setup a test server that returns a simple JSON response
	responseData := map[string]string{"hello": "world"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(responseData)
		w.Write(data)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	client := flyHttp.NewEsiHttpClient(ts.URL, logger, authClient, cache)

	var result map[string]string
	err := client.GetJSON("/", nil, false, &result)
	require.NoError(t, err)
	assert.Equal(t, "world", result["hello"])
}

func TestAPIClient_GetJSON_UsesCache(t *testing.T) {
	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	cachedResponse := map[string]string{"cached": "value"}
	cachedBytes, _ := json.Marshal(cachedResponse)

	// Mock the cache "Get" call to return the cached data
	cache.On("Get", "http://example.com/data").Return(cachedBytes, true)

	client := flyHttp.NewEsiHttpClient("http://example.com", logger, authClient, cache)
	var result map[string]string

	err := client.GetJSON("/data", nil, true, &result)
	require.NoError(t, err)
	assert.Equal(t, "value", result["cached"], "Should use cached data")

	cache.AssertExpectations(t)
	authClient.AssertExpectations(t)
}

func TestAPIClient_GetJSON_TokenRefreshOnUnauthorized(t *testing.T) {
	callCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp := map[string]string{"refreshed": "true"}
		data, _ := json.Marshal(resp)
		w.Write(data)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	// Mock RefreshToken call, which is expected when a 401 is received
	authClient.On("RefreshToken", "refresh-token").
		Return(&oauth2.Token{AccessToken: "new-access-token"}, nil).
		Once()

	client := flyHttp.NewEsiHttpClient(ts.URL, logger, authClient, cache)

	token := &oauth2.Token{
		AccessToken:  "old-access-token",
		RefreshToken: "refresh-token",
	}

	var result map[string]string
	err := client.GetJSON("/", token, false, &result)
	require.NoError(t, err)
	assert.Equal(t, "true", result["refreshed"])

	cache.AssertExpectations(t)
	authClient.AssertExpectations(t)
}

func TestAPIClient_GetJSON_RetryOnServiceUnavailable(t *testing.T) {
	// Setup a server that returns 503 for the first two requests, then succeeds
	callCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		resp := map[string]string{"status": "ok"}
		data, _ := json.Marshal(resp)
		w.Write(data)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	client := flyHttp.NewEsiHttpClient(ts.URL, logger, authClient, cache)

	var result map[string]string
	err := client.GetJSON("/", nil, false, &result)
	require.NoError(t, err)
	assert.Equal(t, 3, callCount, "should have retried twice and succeeded on the third call")
	assert.Equal(t, "ok", result["status"])
}

func TestAPIClient_GetJSON_FailsAfterMaxRetries(t *testing.T) {
	// Setup a server that always returns 503
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	client := flyHttp.NewEsiHttpClient(ts.URL, logger, authClient, cache)

	var result map[string]string
	err := client.GetJSON("/", nil, false, &result)
	require.Error(t, err)

	var cErr *flyErrors.CustomError
	assert.True(t, errors.As(err, &cErr))
	assert.Equal(t, http.StatusServiceUnavailable, cErr.StatusCode, "should return the final error from server")
}

func TestAPIClient_GetJSON_NoToken_NoAuthHeaders(t *testing.T) {
	// Test to ensure no authorization header is sent if no token is provided
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			w.WriteHeader(http.StatusBadRequest) // We expected no auth header
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"success": true}`)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	logger := &testutil.MockLogger{}
	authClient := &testutil.MockAuthClient{}
	cache := &testutil.MockCacheService{}

	client := flyHttp.NewEsiHttpClient(ts.URL, logger, authClient, cache)

	var result map[string]bool
	err := client.GetJSON("/", nil, false, &result)
	require.NoError(t, err)
	assert.True(t, result["success"])
}
