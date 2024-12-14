// auth/auth_client.go
package account

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.AuthClient = (*authClient)(nil)

const (
	tokenURL        = "https://login.eveonline.com/v2/oauth/token"
	requestTimeout  = 10 * time.Second
	contentType     = "application/x-www-form-urlencoded"
	authorization   = "Authorization"
	contentTypeName = "Content-Type"
)

// authClient is a concrete implementation of AuthClient
type authClient struct {
	logger interfaces.Logger
	config *oauth2.Config
	client *http.Client
}

// NewAuthClient initializes and returns an AuthClient implementation.
func NewAuthClient(logger interfaces.Logger, clientID, clientSecret, callbackURL string) interfaces.AuthClient {
	return &authClient{
		logger: logger,
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  callbackURL,
			Scopes: []string{
				"publicData",
				"esi-location.read_location.v1",
				"esi-skills.read_skills.v1",
				"esi-clones.read_clones.v1",
				"esi-clones.read_implants.v1",
				"esi-skills.read_skillqueue.v1",
				"esi-characters.read_corporation_roles.v1",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
				TokenURL: tokenURL,
			},
		},
		// Inject a custom HTTP client if needed, otherwise use default
		client: &http.Client{Timeout: requestTimeout},
	}
}

// GetAuthURL returns the URL for OAuth2 authentication
func (a *authClient) GetAuthURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// ExchangeCode exchanges the authorization code for an access token
func (a *authClient) ExchangeCode(code string) (*oauth2.Token, error) {
	token, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		a.logger.Errorf("Failed to exchange code: %v", err)
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

// RefreshToken performs a token refresh using the current oauth2.Config
func (a *authClient) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		a.logger.Errorf("Failed to create request to refresh token: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(contentTypeName, contentType)
	req.Header.Add(authorization, "Basic "+base64.StdEncoding.EncodeToString([]byte(a.config.ClientID+":"+a.config.ClientSecret)))

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("Failed to make request to refresh token: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			a.logger.Infof("Failed to read response body: %v", readErr)
			return nil, fmt.Errorf("failed to read response body: %w", readErr)
		}
		bodyString := string(bodyBytes)

		a.logger.Warnf("Received non-OK status code %d for request to refresh token. Response body: %s", resp.StatusCode, bodyString)
		return nil, fmt.Errorf("received non-OK status code %d: %s", resp.StatusCode, bodyString)
	}

	var token oauth2.Token
	if decodeErr := json.NewDecoder(resp.Body).Decode(&token); decodeErr != nil {
		a.logger.Errorf("Failed to decode response body: %v", decodeErr)
		return nil, fmt.Errorf("failed to decode response: %w", decodeErr)
	}

	return &token, nil
}
