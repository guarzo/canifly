// esi/helpers.go

package esi

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/auth"
	errors2 "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/persist"
)

const (
	maxRetries = 5
	baseDelay  = 1 * time.Second
	maxDelay   = 32 * time.Second
)

// retryWithExponentialBackoff retries the given function with exponential backoff
func retryWithExponentialBackoff(operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	delay := baseDelay

	for i := 0; i < maxRetries; i++ {
		if result, err = operation(); err == nil {
			return result, nil
		}

		// Retry only on specific error types
		var customErr *errors2.CustomError
		if !errors.As(err, &customErr) || (customErr.StatusCode != http.StatusServiceUnavailable && customErr.StatusCode != http.StatusGatewayTimeout) && customErr.StatusCode != http.StatusInternalServerError {
			break
		}

		if i == maxRetries-1 {
			break
		}

		// Random jitter helps prevent retry storms
		jitter := time.Duration(rand.Int63n(int64(delay)))
		time.Sleep(delay + jitter)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return nil, err
}

func makeRequest(url string, token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if token != nil {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && token != nil {
		// Refresh the token
		newToken, err := auth.RefreshToken(token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		*token = *newToken

		// Retry the request with the new access token
		return makeRequest(url, newToken)
	}

	if customErr, exists := errors2.HttpStatusErrors[resp.StatusCode]; exists {
		return nil, customErr
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors2.NewCustomError(resp.StatusCode, "failed request")
	}

	return bodyBytes, nil
}

// createOperation creates an operation function for the given URL and token
func createOperation(url string, token *oauth2.Token) func() (interface{}, error) {
	return func() (interface{}, error) {
		return makeRequest(url, token)
	}
}

func getResults(address string, token *oauth2.Token) ([]byte, error) {
	operation := createOperation(address, token)

	// Use retryWithExponentialBackoff to handle retries
	result, err := retryWithExponentialBackoff(operation)
	if err != nil {
		return nil, err
	}

	bodyBytes, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to convert result to byte slice")
	}

	return bodyBytes, nil
}

func getResultsWithCache(address string, token *oauth2.Token, dataStore *persist.DataStore) ([]byte, error) {
	// Check the cache first
	if cachedData, found := dataStore.GetFromCache(address); found {
		return cachedData, nil // Return cached data
	}

	bodyBytes, err := getResults(address, token)
	if err != nil {
		return nil, err
	}

	// Store the result in the cache
	dataStore.SetToCache(address, bodyBytes, persist.DefaultExpiration)

	return bodyBytes, nil
}
