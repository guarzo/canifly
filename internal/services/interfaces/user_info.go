package interfaces

import (
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

// UserInfoFetcher provides a minimal interface for fetching user info
type UserInfoFetcher interface {
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
}
