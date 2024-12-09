package handlers

import (
	"fmt"
	"net/http"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	sessionService *flyHttp.SessionService
	dataStore      *persist.DataStore
	logger         *logrus.Logger
}

func NewAccountHandler(s *flyHttp.SessionService, l *logrus.Logger, data *persist.DataStore) *AccountHandler {
	return &AccountHandler{
		sessionService: s,
		logger:         l,
		dataStore:      data,
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

		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			respondError(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		var accountToUpdate *model.Account
		for i, account := range accounts {
			if account.ID == request.AccountID {
				accountToUpdate = &accounts[i]
				break
			}
		}

		if accountToUpdate == nil {
			respondError(w, "Account not found", http.StatusNotFound)
			return
		}

		accountToUpdate.Name = request.AccountName

		if err = h.dataStore.SaveAccounts(accounts); err != nil {
			respondError(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
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
			respondError(w, "AccountID is required", http.StatusBadRequest)
			return
		}

		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			respondError(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		var accountFound bool
		for i := range accounts {
			if accounts[i].ID == request.AccountID {
				if accounts[i].Status == "Alpha" {
					accounts[i].Status = "Omega"
				} else {
					accounts[i].Status = "Alpha"
				}
				accountFound = true
				break
			}
		}

		if !accountFound {
			respondError(w, "Account not found", http.StatusNotFound)
			return
		}

		if err = h.dataStore.SaveAccounts(accounts); err != nil {
			respondError(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
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

		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			respondError(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		index, found := findAccountIndex(accounts, request.AccountName)
		if !found {
			respondError(w, fmt.Sprintf("%s not found in accounts", request.AccountName), http.StatusNotFound)
			return
		}

		// Remove the account
		accounts = append(accounts[:index], accounts[index+1:]...)

		if err := h.dataStore.SaveAccounts(accounts); err != nil {
			respondError(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// helper function to find account by name
func findAccountIndex(accounts []model.Account, accountName string) (int, bool) {
	for i, account := range accounts {
		if account.Name == accountName {
			return i, true
		}
	}
	return 0, false
}
