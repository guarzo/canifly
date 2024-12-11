// services/interfaces/state_service.go
package interfaces

import "github.com/guarzo/canifly/internal/model"

type StateService interface {
	GetAppState() model.AppState
	SetAppStateLogin(isLoggedIn bool) error
	UpdateAndSaveAppState(data model.AppState) error
	ClearAppState()
}

type StateRepository interface {
	SetAppState(as model.AppState)
	GetAppState() model.AppState
	SetAppStateLogin(in bool) error
	ClearAppState()
	SaveAppStateSnapshot(as model.AppState) error
}
