// services/interfaces/dashboard_service.go
package interfaces

import (
	"github.com/guarzo/canifly/internal/model"
)

type DashboardService interface {
	RefreshAccountsAndState() (model.AppState, error)
	RefreshAccounts() ([]model.Account, error)
	PrepareAppData(accounts []model.Account) model.AppState
	UpdateAndSaveAppState(data model.AppState) error
	// If needed, a method to trigger background refresh
	RefreshDataInBackground() error
	// Anything else DashboardHandler needs
}
