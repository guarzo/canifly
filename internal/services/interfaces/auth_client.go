// services/interfaces/auth_client.go
package interfaces

import "golang.org/x/oauth2"

// AuthClient defines methods related to authentication and token management.
type AuthClient interface {
	// RefreshToken takes a refresh token and returns a new, updated oauth2.Token
	RefreshToken(refreshToken string) (*oauth2.Token, error)
}
