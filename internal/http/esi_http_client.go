// internal/services/http/api_client.go
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/persist/eve"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	maxRetries = 5
	baseDelay  = 1 * time.Second
	maxDelay   = 32 * time.Second
)

var _ interfaces.EsiHttpClient = (*EsiHttpClient)(nil)

type EsiHttpClient struct {
	BaseURL      string
	HTTPClient   *http.Client
	Logger       interfaces.Logger
	AuthClient   interfaces.AuthClient
	CacheService interfaces.CacheService
}

// NewEsiHttpClient initializes and returns an EsiHttpClient instance
func NewEsiHttpClient(baseURL string, logger interfaces.Logger, auth interfaces.AuthClient, cache interfaces.CacheService) *EsiHttpClient {
	return &EsiHttpClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger:       logger,
		AuthClient:   auth,
		CacheService: cache,
	}
}

// GetJSON retrieves JSON data from the specified endpoint. It supports optional caching and token usage.
// If `useCache` is true, it will attempt to return cached data before making a request.
// If a token is provided, it will include it in the request and attempt token refresh if Unauthorized.
func (c *EsiHttpClient) GetJSON(endpoint string, token *oauth2.Token, useCache bool, target interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	return c.GetJSONFromURL(url, token, useCache, target)
}

func (c *EsiHttpClient) GetJSONFromURL(url string, token *oauth2.Token, useCache bool, target interface{}) error {
	// Check cache first
	if useCache && c.CacheService != nil {
		if cachedData, found := c.CacheService.Get(url); found {
			c.Logger.Debugf("using cached data for %s", url)
			return json.Unmarshal(cachedData, target)
		} else {
			c.Logger.Debugf("no cached data found for %s", url)
		}
	}

	// Define the operation for retry
	operation := func() ([]byte, error) {
		return c.doRequestWithToken("GET", url, nil, token)
	}

	bodyBytes, err := c.retryWithExponentialBackoff(operation)
	if err != nil {
		return err
	}

	// Cache the response if needed
	if useCache && c.CacheService != nil {
		c.CacheService.Set(url, bodyBytes, eve.DefaultExpiration)
	}

	return json.Unmarshal(bodyBytes, target)
}

// doRequestWithToken performs a request and handles token refresh if necessary.
func (c *EsiHttpClient) doRequestWithToken(method, url string, body interface{}, token *oauth2.Token) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			c.Logger.WithError(err).Error("Failed to serialize request body")
			return nil, fmt.Errorf("failed to serialize request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil && token.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to execute request")
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) && token != nil && token.RefreshToken != "" {
		// Attempt token refresh
		newToken, refreshErr := c.AuthClient.RefreshToken(token.RefreshToken)
		if refreshErr != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", refreshErr)
		}
		token.AccessToken = newToken.AccessToken
		// Retry once with the new token
		return c.doRequestWithToken(method, url, body, token)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("Received non-2xx response")

		return nil, flyErrors.NewCustomError(resp.StatusCode, fmt.Sprintf("unexpected status code: %d, response: %s", resp.StatusCode, body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

// retryWithExponentialBackoff attempts the given operation multiple times with exponential backoff on certain HTTP errors.
func (c *EsiHttpClient) retryWithExponentialBackoff(operation func() ([]byte, error)) ([]byte, error) {
	delay := baseDelay
	for i := 0; i < maxRetries; i++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}

		var customErr *flyErrors.CustomError
		if !shouldRetry(err, &customErr) {
			return nil, err
		}

		// If we're at the last attempt, don't retry
		if i == maxRetries-1 {
			return nil, err
		}

		jitter := time.Duration(rand.Int63n(int64(delay)))
		time.Sleep(delay + jitter)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return nil, fmt.Errorf("exceeded maximum retries")
}

// shouldRetry checks if the error status code warrants a retry.
func shouldRetry(err error, customErr **flyErrors.CustomError) bool {
	if errors.As(err, customErr) {
		switch (*customErr).StatusCode {
		case http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusInternalServerError:
			return true
		}
	}
	return false
}
