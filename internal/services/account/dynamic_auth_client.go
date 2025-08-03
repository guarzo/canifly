package account

import (
	"golang.org/x/oauth2"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// DynamicAuthClient wraps the regular auth client but loads credentials dynamically
type DynamicAuthClient struct {
	logger         interfaces.Logger
	configService  interfaces.ConfigurationService
	baseCallbackURL string
}

func NewDynamicAuthClient(logger interfaces.Logger, configService interfaces.ConfigurationService, baseCallbackURL string) interfaces.AuthClient {
	return &DynamicAuthClient{
		logger:          logger,
		configService:   configService,
		baseCallbackURL: baseCallbackURL,
	}
}

func (d *DynamicAuthClient) getAuthClient() (interfaces.AuthClient, error) {
	// Get credentials from config service
	clientID, clientSecret, callbackURL, err := d.configService.GetEVECredentials()
	if err != nil {
		d.logger.Errorf("Failed to get EVE credentials: %v", err)
		return nil, err
	}
	
	d.logger.Infof("Dynamic auth client loaded credentials - ClientID: %s, CallbackURL: %s", clientID, callbackURL)
	
	// Use stored callback URL if available, otherwise use base
	if callbackURL == "" {
		callbackURL = d.baseCallbackURL
	}
	
	// Create a new auth client with current credentials
	return NewAuthClient(d.logger, clientID, clientSecret, callbackURL), nil
}

func (d *DynamicAuthClient) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	client, err := d.getAuthClient()
	if err != nil {
		return nil, err
	}
	return client.RefreshToken(refreshToken)
}

func (d *DynamicAuthClient) GetAuthURL(state string) string {
	client, err := d.getAuthClient()
	if err != nil {
		d.logger.Errorf("Failed to get auth client: %v", err)
		return ""
	}
	return client.GetAuthURL(state)
}

func (d *DynamicAuthClient) ExchangeCode(code string) (*oauth2.Token, error) {
	client, err := d.getAuthClient()
	if err != nil {
		return nil, err
	}
	return client.ExchangeCode(code)
}