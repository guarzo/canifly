package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

func UpdateAccountNameHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			AccountID   int64  `json:"accountID"`
			AccountName string `json:"accountName"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate the new account name
		if request.AccountName == "" {
			http.Error(w, "Account name cannot be empty", http.StatusBadRequest)
			return
		}

		// Retrieve the session
		session, err := s.Get(r, sessionName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
			return
		}

		// Retrieve the logged-in user (main identity)
		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			http.Error(w, "Main identity not found in session", http.StatusUnauthorized)
			return
		}

		// Fetch accounts for the logged-in user
		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the account to update
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

		// Update the account name
		accountToUpdate.Name = request.AccountName

		// Save the updated accounts
		err = persist.SaveAccounts(accounts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func ToggleAccountStatusHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			AccountID int64 `json:"accountID"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if request.AccountID == 0 {
			http.Error(w, "AccountID is required", http.StatusBadRequest)
			return
		}

		// Retrieve the session
		session, err := s.Get(r, sessionName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
			return
		}

		// Retrieve the logged-in user
		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Fetch accounts for the user
		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the account and toggle its status
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

		// Save the updated accounts
		err = persist.SaveAccounts(accounts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
