// handlers/account_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AccountHandler struct {
	sessionService interfaces.SessionService
	accountService interfaces.AccountManagementService
	logger         interfaces.Logger
}

func NewAccountHandler(session interfaces.SessionService, logger interfaces.Logger, accountSrv interfaces.AccountManagementService) *AccountHandler {
	return &AccountHandler{
		sessionService: session,
		logger:         logger,
		accountService: accountSrv,
	}
}

// RESTful endpoint: GET /api/accounts
func (h *AccountHandler) ListAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := h.accountService.FetchAccounts()
		if err != nil {
			respondError(w, "Failed to fetch accounts", http.StatusInternalServerError)
			return
		}
		
		respondJSON(w, accounts)
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

		accounts, err := h.accountService.FetchAccounts()
		if err != nil {
			respondError(w, "Failed to fetch accounts", http.StatusInternalServerError)
			return
		}

		for _, acc := range accounts {
			if acc.ID == accountID {
				respondJSON(w, acc)
				return
			}
		}

		respondError(w, "Account not found", http.StatusNotFound)
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
			Name       *string `json:"name,omitempty"`
			IsActive   *bool   `json:"isActive,omitempty"`
			IsVisible  *bool   `json:"isVisible,omitempty"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Update name if provided
		if request.Name != nil {
			if *request.Name == "" {
				respondError(w, "Account name cannot be empty", http.StatusBadRequest)
				return
			}
			err := h.accountService.UpdateAccountName(accountID, *request.Name)
			if err != nil {
				respondError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Toggle active status if provided
		if request.IsActive != nil {
			err := h.accountService.ToggleAccountStatus(accountID)
			if err != nil {
				if err.Error() == "account not found" {
					respondError(w, "Account not found", http.StatusNotFound)
				} else {
					respondError(w, fmt.Sprintf("Failed to toggle account status: %v", err), http.StatusInternalServerError)
				}
				return
			}
		}

		// Toggle visibility if provided
		if request.IsVisible != nil {
			err := h.accountService.ToggleAccountVisibility(accountID)
			if err != nil {
				if err.Error() == "account not found" {
					respondError(w, "Account not found", http.StatusNotFound)
				} else {
					respondError(w, fmt.Sprintf("Failed to toggle account visibility: %v", err), http.StatusInternalServerError)
				}
				return
			}
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// Legacy endpoint handlers (to be deprecated)
func (h *AccountHandler) UpdateAccountName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AccountID   int64  `json:"accountID"`
			AccountName string `json:"accountName"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.AccountName == "" {
			respondError(w, "Account name cannot be empty", http.StatusBadRequest)
			return
		}

		err := h.accountService.UpdateAccountName(request.AccountID, request.AccountName)
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

func (h *AccountHandler) ToggleAccountStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AccountID int64 `json:"accountID"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.AccountID == 0 {
			respondError(w, "UserId is required", http.StatusBadRequest)
			return
		}

		err := h.accountService.ToggleAccountStatus(request.AccountID)
		if err != nil {
			if err.Error() == "account not found" {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to toggle account status: %v", err), http.StatusInternalServerError)
			}
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

func (h *AccountHandler) ToggleAccountVisibility() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AccountID int64 `json:"accountID"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.AccountID == 0 {
			respondError(w, "UserId is required", http.StatusBadRequest)
			return
		}

		err := h.accountService.ToggleAccountVisibility(request.AccountID)
		if err != nil {
			if err.Error() == "account not found" {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to toggle account visbility: %v", err), http.StatusInternalServerError)
			}
			return
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

		// Get account to find its name (needed for RemoveAccountByName)
		accounts, err := h.accountService.FetchAccounts()
		if err != nil {
			respondError(w, "Failed to fetch accounts", http.StatusInternalServerError)
			return
		}

		var accountName string
		for _, acc := range accounts {
			if acc.ID == accountID {
				accountName = acc.Name
				break
			}
		}

		if accountName == "" {
			respondError(w, "Account not found", http.StatusNotFound)
			return
		}

		err = h.accountService.RemoveAccountByName(accountName)
		if err != nil {
			if err.Error() == fmt.Sprintf("account %s not found", accountName) {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to remove account: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Legacy endpoint (to be deprecated)
func (h *AccountHandler) RemoveAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AccountName string `json:"accountName"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if request.AccountName == "" {
			respondError(w, "AccountName is required", http.StatusBadRequest)
			return
		}

		err := h.accountService.RemoveAccountByName(request.AccountName)
		if err != nil {
			if err.Error() == fmt.Sprintf("account %s not found", request.AccountName) {
				respondError(w, "Account not found", http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to remove account: %v", err), http.StatusInternalServerError)
			}
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}
