package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type DashboardHandler struct {
	sessionService   interfaces.SessionService
	logger           interfaces.Logger
	dashboardService interfaces.DashboardService
}

func NewDashboardHandler(
	s interfaces.SessionService,
	logger interfaces.Logger,
	dashboardService interfaces.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		sessionService:   s,
		logger:           logger,
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) handleAppStateRefresh(w http.ResponseWriter, noCache bool) {
	// Get current app state (cached) from the DashboardService (which uses AppStateService internally)
	appState := h.dashboardService.GetCurrentAppState()

	if !noCache && len(appState.AccountData.Accounts) > 0 {
		// We have cached data
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(appState)
		// Trigger background refresh
		go func() {
			if err := h.dashboardService.RefreshDataInBackground(); err != nil {
				h.logger.Errorf("background refresh failed: %v", err)
			}
		}()
		return
	}

	// If noCache is true or we have no cached accounts, refresh and rebuild state
	updatedData, err := h.dashboardService.RefreshAccountsAndState()
	if err != nil {
		h.logger.Errorf("Failed to validate accounts: %v", err)
		http.Error(w, `{"error":"Failed to validate accounts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedData); err != nil {
		http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
		return
	}
}

func (h *DashboardHandler) GetDashboardData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infof("GetAppData Called")
		h.handleAppStateRefresh(w, false)
	}
}

func (h *DashboardHandler) GetDashboardDataNoCache() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infof("GetAppDataNoCache Called")
		h.handleAppStateRefresh(w, true)
	}
}
