package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPClient interface that can be used by the esi package
type HTTPClient interface {
	DoRequest(method, endpoint string, body interface{}, target interface{}) error
}

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *logrus.Logger
	Token      string // OAuth token or API key
}

// NewAPIClient initializes and returns an APIClient instance
func NewAPIClient(baseURL, token string, logger *logrus.Logger) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, // Set a timeout for API requests
		},
		Logger: logger,
		Token:  token,
	}
}

// DoRequest sends an HTTP request and parses the response
func (c *APIClient) DoRequest(method, endpoint string, body interface{}, target interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			c.Logger.WithError(err).Error("Failed to serialize request body")
			return fmt.Errorf("failed to serialize request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to create request")
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to execute request")
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("Received non-2xx response")
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to read response body")
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if target != nil {
		if err := json.Unmarshal(respBody, target); err != nil {
			c.Logger.WithError(err).Error("Failed to parse JSON response")
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
	}

	return nil
}
