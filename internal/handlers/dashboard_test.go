package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/guarzo/canifly/internal/handlers"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardService is a testify mock for DashboardService
type MockDashboardService struct {
	mock.Mock
	callCount int64
	wg        sync.WaitGroup
}

func (m *MockDashboardService) RefreshAccountsAndState() (model.AppState, error) {
	args := m.Called()
	return args.Get(0).(model.AppState), args.Error(1)
}

func (m *MockDashboardService) RefreshDataInBackground() error {
	defer m.wg.Done() // signal that this call was made
	args := m.Called()
	return args.Error(0)
}

func (m *MockDashboardService) GetCurrentAppState() model.AppState {
	args := m.Called()
	return args.Get(0).(model.AppState)
}

// Helper to set lastRefreshTime in the handler for tests.
func setLastRefreshTimeForTest(h *handlers.DashboardHandler, t time.Time) {
	// Using reflection or a test helper within the same package is cleaner.
	// If handlers package is different, you can either export the field or add a helper method in production code under a test build tag.
	h.SetLastRefreshTimeForTest(t)
}

func TestDashboardHandler_CachedDataWithBackgroundRefresh(t *testing.T) {
	logger := &testutil.MockLogger{}
	sessionSvc := &testutil.MockSessionService{}
	dashboardSvc := &MockDashboardService{}
	dashboardSvc.wg.Add(1) // Expect one background refresh call
	handler := handlers.NewDashboardHandler(sessionSvc, logger, dashboardSvc)

	cachedState := model.AppState{
		AccountData: model.AccountData{
			Accounts: []model.Account{{Name: "TestAcc"}},
		},
	}
	dashboardSvc.On("GetCurrentAppState").Return(cachedState)
	dashboardSvc.On("RefreshDataInBackground").Return(nil).Once()

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	rr := httptest.NewRecorder()
	handler.GetDashboardData().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.AppState
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, cachedState, resp)

	// Wait for the goroutine to run RefreshDataInBackground
	dashboardSvc.wg.Wait()

	dashboardSvc.AssertExpectations(t)
}

func TestDashboardHandler_CachedDataSkipBackgroundRefresh(t *testing.T) {
	logger := &testutil.MockLogger{}
	sessionSvc := &testutil.MockSessionService{}
	dashboardSvc := &MockDashboardService{}
	handler := handlers.NewDashboardHandler(sessionSvc, logger, dashboardSvc)

	// Set lastRefreshTime to now, so background refresh should be skipped
	setLastRefreshTimeForTest(handler, time.Now())

	cachedState := model.AppState{
		AccountData: model.AccountData{
			Accounts: []model.Account{{Name: "TestAcc2"}},
		},
	}
	dashboardSvc.On("GetCurrentAppState").Return(cachedState)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	rr := httptest.NewRecorder()
	handler.GetDashboardData().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.AppState
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, cachedState, resp)

	dashboardSvc.AssertNotCalled(t, "RefreshDataInBackground")
}

func TestDashboardHandler_NoCacheRefresh(t *testing.T) {
	logger := &testutil.MockLogger{}
	sessionSvc := &testutil.MockSessionService{}
	dashboardSvc := &MockDashboardService{}
	handler := handlers.NewDashboardHandler(sessionSvc, logger, dashboardSvc)

	freshState := model.AppState{
		AccountData: model.AccountData{
			Accounts: []model.Account{{Name: "RefreshedAcc"}},
		},
	}
	dashboardSvc.On("GetCurrentAppState").Return(model.AppState{})
	dashboardSvc.On("RefreshAccountsAndState").Return(freshState, nil).Once()

	req, _ := http.NewRequest("GET", "/dashboard/nocache", nil)
	rr := httptest.NewRecorder()
	handler.GetDashboardDataNoCache().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.AppState
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, freshState, resp)

	dashboardSvc.AssertExpectations(t)
}

func TestDashboardHandler_NoCacheRefreshError(t *testing.T) {
	logger := &testutil.MockLogger{}
	sessionSvc := &testutil.MockSessionService{}
	dashboardSvc := &MockDashboardService{}
	handler := handlers.NewDashboardHandler(sessionSvc, logger, dashboardSvc)

	dashboardSvc.On("GetCurrentAppState").Return(model.AppState{})
	dashboardSvc.On("RefreshAccountsAndState").Return(model.AppState{}, assert.AnError).Once()

	req, _ := http.NewRequest("GET", "/dashboard/nocache", nil)
	rr := httptest.NewRecorder()
	handler.GetDashboardDataNoCache().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to validate accounts")

	dashboardSvc.AssertExpectations(t)
}

func TestDashboardHandler_BadJSONResponse(t *testing.T) {
	logger := &testutil.MockLogger{}
	sessionSvc := &testutil.MockSessionService{}
	dashboardSvc := &MockDashboardService{}
	handler := handlers.NewDashboardHandler(sessionSvc, logger, dashboardSvc)

	freshState := model.AppState{
		AccountData: model.AccountData{
			Accounts: []model.Account{{Name: "AccForBadJSON"}},
		},
	}
	dashboardSvc.On("GetCurrentAppState").Return(model.AppState{})
	dashboardSvc.On("RefreshAccountsAndState").Return(freshState, nil).Once()

	req, _ := http.NewRequest("GET", "/dashboard/nocache", nil)

	// Wrap the recorder in a writer that fails only on the first write
	rr := httptest.NewRecorder()
	fw := &failingResponseWriter{ResponseRecorder: rr, failOnce: true}

	handler.GetDashboardDataNoCache().ServeHTTP(fw, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	// Now we can check the body for the error message, since the second write should succeed
	assert.Contains(t, rr.Body.String(), "Failed to encode data")

	dashboardSvc.AssertExpectations(t)
}

// failingResponseWriter fails only on the first write attempt, then succeeds.
type failingResponseWriter struct {
	*httptest.ResponseRecorder
	failOnce bool
}

func (f *failingResponseWriter) Write(p []byte) (int, error) {
	if f.failOnce {
		f.failOnce = false
		return 0, assert.AnError
	}
	return f.ResponseRecorder.Write(p)
}
