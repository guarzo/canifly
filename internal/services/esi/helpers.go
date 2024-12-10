// services/esi/helpers.go
package esi

import (
	"errors"
	"fmt"
	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/sirupsen/logrus"
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

func makeRequest(url string, token *oauth2.Token, auth auth.AuthClient) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}

	if token != nil && token.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && token != nil {
		// Refresh the token using auth client
		newToken, err := auth.RefreshToken(token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		*token = *newToken
		// Retry the request with the new access token
		return makeRequest(url, newToken, auth)
	}

	if customErr, exists := flyErrors.HttpStatusErrors[resp.StatusCode]; exists {
		return nil, customErr
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, flyErrors.NewCustomError(resp.StatusCode, "failed request")
	}

	return bodyBytes, nil
}

func createOperation(url string, token *oauth2.Token, auth auth.AuthClient) func() (interface{}, error) {
	return func() (interface{}, error) {
		return makeRequest(url, token, auth)
	}
}

func getResults(address string, token *oauth2.Token, auth auth.AuthClient) ([]byte, error) {
	operation := createOperation(address, token, auth)

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

func getResultsWithCache(address string, token *oauth2.Token, cacheService interfaces.CacheService, logger *logrus.Logger, auth auth.AuthClient) ([]byte, error) {
	if cachedData, found := cacheService.Get(address); found {
		logger.Debugf("using cacheService data for call to %s", address)
		return cachedData, nil
	} else {
		logger.Debugf("no cacheService data found for call to %s", address)
	}

	bodyBytes, err := getResults(address, token, auth)
	if err != nil {
		return nil, err
	}

	cacheService.Set(address, bodyBytes, persist.DefaultExpiration)
	return bodyBytes, nil
}
