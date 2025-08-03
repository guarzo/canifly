// handlers/account_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AccountHandler struct {
	sessionService interfaces.SessionService
	accountService interfaces.AccountManagementService
	logger         interfaces.Logger
	cache          interfaces.HTTPCacheService
	wsHub          *WebSocketHub
}

func NewAccountHandler(session interfaces.SessionService, logger interfaces.Logger, accountSrv interfaces.AccountManagementService, cache interfaces.HTTPCacheService, wsHub *WebSocketHub) *AccountHandler {
	return &AccountHandler{
		sessionService: session,
		logger:         logger,
		accountService: accountSrv,
		cache:          cache,
		wsHub:          wsHub,
	}
}

// RESTful endpoint: GET /api/accounts
func (h *AccountHandler) ListAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse pagination parameters
		paginationParams := ParsePaginationParams(r)

		// Check cache first (cache key includes pagination params)
		cacheKey := fmt.Sprintf("accounts:list:page:%d:limit:%d", paginationParams.Page, paginationParams.Limit)

		cacheHandler := WithCache(
			h.cache,
			h.logger,
			cacheKey,
			5*time.Minute, // Cache for 5 minutes
			func() (interface{}, error) {
				accounts, err := h.accountService.FetchAccounts()
				if err != nil {
					return nil, err
				}

				// Apply pagination to accounts
				paginatedResponse := PaginateAccounts(accounts, paginationParams)

				return paginatedResponse, nil
			},
		)
		cacheHandler(w, r)
	}
}

// RESTful endpoint: GET /api/accounts/:id
func (h *AccountHandler) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		account, err := h.accountService.GetAccountByID(accountID)
		if err != nil {
			if err.Error() == fmt.Sprintf("account with ID %d not found", accountID) {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, "Failed to fetch account", http.StatusInternalServerError)
			}
			return
		}

		respondJSON(w, account)
	}
}

// RESTful endpoint: PATCH /api/accounts/:id
func (h *AccountHandler) UpdateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		var request struct {
			Name      *string `json:"name,omitempty"`
			IsActive  *bool   `json:"isActive,omitempty"`
			IsVisible *bool   `json:"isVisible,omitempty"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Build update request
		updates := interfaces.AccountUpdateRequest{}

		// Validate and set name if provided
		if request.Name != nil {
			if *request.Name == "" {
				respondError(w, "Account name cannot be empty", http.StatusBadRequest)
				return
			}
			updates.Name = request.Name
		}

		// Set status if provided (convert bool to AccountStatus)
		if request.IsActive != nil {
			if *request.IsActive {
				status := model.Omega
				updates.Status = &status
			} else {
				status := model.Alpha
				updates.Status = &status
			}
		}

		// Set visibility if provided
		if request.IsVisible != nil {
			updates.Visible = request.IsVisible
		}

		// Perform atomic update
		err = h.accountService.UpdateAccount(accountID, updates)
		if err != nil {
			if err.Error() == "account not found" {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to update account: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Invalidate accounts cache after successful update
		InvalidateCache(h.cache, "accounts:")

		// Broadcast update via WebSocket
		if h.wsHub != nil {
			h.wsHub.BroadcastUpdate("account:updated", map[string]interface{}{
				"accountId": accountID,
			})
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// RESTful endpoint: DELETE /api/accounts/:id
func (h *AccountHandler) DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		err = h.accountService.RemoveAccountByID(accountID)
		if err != nil {
			if err.Error() == fmt.Sprintf("account with ID %d not found", accountID) {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to remove account: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Invalidate accounts cache after successful deletion
		InvalidateCache(h.cache, "accounts:")

		// Broadcast deletion via WebSocket
		if h.wsHub != nil {
			h.wsHub.BroadcastUpdate("account:deleted", map[string]interface{}{
				"accountId": accountID,
			})
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
