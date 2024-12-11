// services/interfaces/state_service.go
package interfaces

import "github.com/guarzo/canifly/internal/model"

type StateService interface {
	GetAppState() model.AppState
	SetAppStateLogin(isLoggedIn bool) error
	UpdateAndSaveAppState(data model.AppState) error
	ClearAppState()
}
