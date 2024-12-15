package handlers

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type DashboardHandler struct {
	sessionService   interfaces.SessionService
	logger           interfaces.Logger
	dashboardService interfaces.DashboardService
	// lastRefreshTime holds a Unix nano timestamp of the last time a refresh occurred
	lastRefreshTime int64
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
		// lastRefreshTime is initially 0 indicating never refreshed
	}
}

func (h *DashboardHandler) handleAppStateRefresh(w http.ResponseWriter, noCache bool) {
	appState := h.dashboardService.GetCurrentAppState()

	// If we have cached data, and we are allowed to use it (noCache == false):
	if !noCache && len(appState.AccountData.Accounts) > 0 {
		respondEncodedData(w, appState)

		// Attempt a background refresh if it's been more than 5s since the last refresh
		now := time.Now().UnixNano()
		old := atomic.LoadInt64(&h.lastRefreshTime)
		if now-old > 5*int64(time.Second) {
			// Try to set the last refresh time optimistically. If CAS fails, someone else just did a refresh.
			if atomic.CompareAndSwapInt64(&h.lastRefreshTime, old, now) {
				go func() {
					if err := h.dashboardService.RefreshDataInBackground(); err != nil {
						h.logger.Errorf("background refresh failed: %v", err)
					} else {
						h.logger.Debug("Background refresh completed successfully.")
						// Update lastRefreshTime again after successful refresh.
						atomic.StoreInt64(&h.lastRefreshTime, time.Now().UnixNano())
					}
				}()
			} else {
				h.logger.Debug("Another refresh is already in progress, skipping.")
			}
		} else {
			h.logger.Debugf("Skipping background refresh; only %v since last refresh", time.Duration(now-old))
		}
		return
	}

	// If noCache is true or we have no cached accounts, do a full refresh now (ignoring the last refresh time)
	updatedData, err := h.dashboardService.RefreshAccountsAndState()
	if err != nil {
		h.logger.Errorf("Failed to validate accounts: %v", err)
		respondError(w, "Failed to validate accounts", http.StatusInternalServerError)
		return
	}

	// After a successful noCache (or initial) refresh, update the lastRefreshTime
	atomic.StoreInt64(&h.lastRefreshTime, time.Now().UnixNano())

	//w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(updatedData); err != nil {
	//	http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
	//	return
	//}
	respondEncodedData(w, updatedData)
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

// Test helper method
func (h *DashboardHandler) SetLastRefreshTimeForTest(t time.Time) {
	atomic.StoreInt64(&h.lastRefreshTime, t.UnixNano())
}
