package auth

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

	"github.com/guarzo/canifly/internal/utils/xlog"
)

const (
	tokenURL        = "https://login.eveonline.com/v2/oauth/token"
	requestTimeout  = 10 * time.Second
	contentType     = "application/x-www-form-urlencoded"
	authorization   = "Authorization"
	contentTypeName = "Content-Type"
)

var httpClient = &http.Client{
	Timeout: requestTimeout,
}

// AuthClient represents what we need from the auth subsystem.
type AuthClient interface {
	RefreshToken(refreshToken string) (*oauth2.Token, error)
}

var oauth2Config *oauth2.Config

// InitializeOAuth initializes the OAuth2 configuration
func InitializeOAuth(clientID, clientSecret, callbackURL string) {
	oauth2Config = &oauth2.Config{
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
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
	}
}

// GetAuthURL returns the URL for OAuth2 authentication
func GetAuthURL(state string) string {
	return oauth2Config.AuthCodeURL(state)
}

// ExchangeCode exchanges the authorization code for an access token
func ExchangeCode(code string) (*oauth2.Token, error) {
	return oauth2Config.Exchange(context.Background(), code)
}

// RefreshToken performs a token refresh using the global oauth2Config
func RefreshToken(refreshToken string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		xlog.Logf("Failed to create request to refresh token: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(contentTypeName, contentType)
	req.Header.Add(authorization, "Basic "+base64.StdEncoding.EncodeToString([]byte(oauth2Config.ClientID+":"+oauth2Config.ClientSecret)))

	resp, err := httpClient.Do(req)
	if err != nil {
		xlog.Logf("Failed to make request to refresh token: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			xlog.Logf("Failed to read response body: %v", readErr)
			return nil, fmt.Errorf("failed to read response body: %w", readErr)
		}
		bodyString := string(bodyBytes)

		xlog.Logf("Received non-OK status code %d for request to refresh token. Response body: %s", resp.StatusCode, bodyString)
		return nil, fmt.Errorf("received non-OK status code %d: %s", resp.StatusCode, bodyString)
	}

	var token oauth2.Token
	if decodeErr := json.NewDecoder(resp.Body).Decode(&token); decodeErr != nil {
		xlog.Logf("Failed to decode response body: %v", decodeErr)
		return nil, fmt.Errorf("failed to decode response: %w", decodeErr)
	}

	return &token, nil
}

// authClient is a concrete implementation of AuthClient
type authClient struct{}

// NewAuthClient returns a concrete AuthClient implementation
func NewAuthClient() AuthClient {
	return &authClient{}
}

func (a *authClient) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return RefreshToken(refreshToken)
}
