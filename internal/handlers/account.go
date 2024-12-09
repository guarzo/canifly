package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

type AccountHandler struct {
	sessionService *http2.SessionService
	dataStore      *persist.DataStore
	logger         *logrus.Logger
}

func NewAccountHandler(s *http2.SessionService, l *logrus.Logger, data *persist.DataStore) *AccountHandler {
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
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.AccountName == "" {
			http.Error(w, "Account name cannot be empty", http.StatusBadRequest)
			return
		}

		session, err := h.sessionService.Get(r, http2.SessionName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
			return
		}

		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
		if !ok || loggedIn == 0 {
			http.Error(w, "Main identity not found in session", http.StatusUnauthorized)
			return
		}

		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
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
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}

		accountToUpdate.Name = request.AccountName

		if err = h.dataStore.SaveAccounts(accounts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func (h *AccountHandler) ToggleAccountStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AccountID int64 `json:"accountID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.AccountID == 0 {
			http.Error(w, "AccountID is required", http.StatusBadRequest)
			return
		}

		session, err := h.sessionService.Get(r, http2.SessionName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
			return
		}

		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
		if !ok || loggedIn == 0 {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
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
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}

		if err = h.dataStore.SaveAccounts(accounts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
