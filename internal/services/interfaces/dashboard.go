// services/interfaces/dashboard_service.go
package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
)

type DashboardService interface {
	RefreshAccountsAndState() (model.AppState, error)
	RefreshDataInBackground() error
	GetCurrentAppState() model.AppState
}
