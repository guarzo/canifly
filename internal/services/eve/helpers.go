// services/esi/helpers.go
package eve

import (
	"encoding/json"
	"errors"
	"fmt"
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

func retryWithExponentialBackoff(operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	delay := baseDelay

	for i := 0; i < maxRetries; i++ {
		if result, err = operation(); err == nil {
			return result, nil
		}

		var customErr *flyErrors.CustomError
		if !errors.As(err, &customErr) || (customErr.StatusCode != http.StatusServiceUnavailable && customErr.StatusCode != http.StatusGatewayTimeout && customErr.StatusCode != http.StatusInternalServerError) {
			break
		}

		if i == maxRetries-1 {
			break
		}

		jitter := time.Duration(rand.Int63n(int64(delay)))
		time.Sleep(delay + jitter)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return nil, err
}

func makeRequest(url string, token *oauth2.Token, auth interfaces.AuthClient, httpClient interfaces.HTTPClient) ([]byte, error) {
	var raw json.RawMessage
	err := httpClient.DoRequest("GET", url, nil, &raw)
	if err != nil {
		var customErr *flyErrors.CustomError
		if errors.As(err, &customErr) && customErr.StatusCode == http.StatusUnauthorized && token != nil && token.RefreshToken != "" {
			// Refresh the token
			newToken, refreshErr := auth.RefreshToken(token.RefreshToken)
			if refreshErr != nil {
				return nil, fmt.Errorf("failed to refresh token: %w", refreshErr)
			}
			*token = *newToken
			// Retry the request with the new token
			err = httpClient.DoRequest("GET", url, nil, &raw)
			if err != nil {
				return nil, err // If still fails, return the error
			}
		} else {
			// If not 401 or can't refresh, return the original error
			return nil, err
		}
	}

	// Successful response (2xx)
	return raw, nil
}

func createOperation(url string, token *oauth2.Token, auth interfaces.AuthClient, client interfaces.HTTPClient) func() (interface{}, error) {
	return func() (interface{}, error) {
		return makeRequest(url, token, auth, client)
	}
}

func getResults(address string, token *oauth2.Token, auth interfaces.AuthClient, client interfaces.HTTPClient) ([]byte, error) {
	operation := createOperation(address, token, auth, client)

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

func getResultsWithCache(address string, token *oauth2.Token, cacheService interfaces.CacheService, logger interfaces.Logger, auth interfaces.AuthClient, client interfaces.HTTPClient) ([]byte, error) {
	if cachedData, found := cacheService.Get(address); found {
		logger.Debugf("using cacheService data for call to %s", address)
		return cachedData, nil
	} else {
		logger.Debugf("no cacheService data found for call to %s", address)
	}

	bodyBytes, err := getResults(address, token, auth, client)
	if err != nil {
		return nil, err
	}

	cacheService.Set(address, bodyBytes, eve.DefaultExpiration)
	return bodyBytes, nil
}
