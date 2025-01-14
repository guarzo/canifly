// handlers/account_handler.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AccountHandler struct {
	sessionService interfaces.SessionService
	accountService interfaces.AccountService
	logger         interfaces.Logger
}

func NewAccountHandler(session interfaces.SessionService, logger interfaces.Logger, accountSrv interfaces.AccountService) *AccountHandler {
	return &AccountHandler{
		sessionService: session,
		logger:         logger,
		accountService: accountSrv,
	}
}

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
