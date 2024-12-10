// services/interfaces/state_service.go
package interfaces

import "github.com/guarzo/canifly/internal/model"

type StateService interface {
	GetAppState() model.AppState
	SetAppState(appState model.AppState) error
	SaveAppStateSnapshot(appState model.AppState) error
	UpdateAndSaveAppState(appState model.AppState) error
}
