package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
)

// LoginRepository interface is used for temporary OAuth state management
// This repository stores in-memory state for the OAuth login flow
type LoginRepository interface {
	Set(state string, authStatus *model.AuthStatus)
	Get(state string) (*model.AuthStatus, bool)
	Delete(state string)
}