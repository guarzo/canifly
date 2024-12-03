package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func UpdateCharacterHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			CharacterID int64                  `json:"characterID"`
			Updates     map[string]interface{} `json:"updates"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate input
		if request.CharacterID == 0 || len(request.Updates) == 0 {
			http.Error(w, "CharacterID and updates are required", http.StatusBadRequest)
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
		accounts, err := persist.FetchAccountByIdentity(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the character and update it
		var characterFound bool
		for i := range accounts {
			for j := range accounts[i].Characters {
				charIdentity := &accounts[i].Characters[j]
				if charIdentity.Character.CharacterID == request.CharacterID {
					// Update character properties
					for key, value := range request.Updates {
						switch key {
						case "Role":
							if roleStr, ok := value.(string); ok {
								charIdentity.Role = roleStr
							}
						case "MCT":
							if mctBool, ok := value.(bool); ok {
								charIdentity.MCT = mctBool
							}
						default:
							http.Error(w, fmt.Sprintf("Unknown update field: %s", key), http.StatusBadRequest)
							return
						}
					}
					characterFound = true
					break
				}
			}
			if characterFound {
				break
			}
		}

		if !characterFound {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}

		// Save the updated accounts
		err = persist.SaveAccounts(mainIdentity, accounts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func GetUnassignedCharactersHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xlog.Logf("start of get unassigned characters handler")

		// Start or retrieve the session
		session, _ := s.Get(r, sessionName)

		// Retrieve the main identity from the session
		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, "Attempt to get unassigned characters without main identity", "/logout")
			return
		}

		// Fetch the unassigned characters for the logged-in user
		unassignedCharacters, err := persist.FetchUnassignedCharacters(mainIdentity)
		if err != nil {
			http.Error(w, "Error fetching unassigned characters", http.StatusInternalServerError)
			return
		}
		xlog.Logf("%v", unassignedCharacters)
		// Return the unassigned characters in the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(unassignedCharacters)
		if err != nil {
			http.Error(w, "Error encoding unassigned characters", http.StatusInternalServerError)
		}
	}
}

func AssignCharacterHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the incoming request
		var request struct {
			CharacterID int64  `json:"characterID"`
			AccountID   *int64 `json:"accountID"` // Optional; if nil, create a new account
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
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

		// Fetch accounts and unassigned characters
		accounts, err := persist.FetchAccountByIdentity(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the character to assign
		var characterToAssign *model.CharacterIdentity
		for i, character := range unassignedCharacters {
			if character.Character.CharacterID == request.CharacterID {
				characterToAssign = &unassignedCharacters[i]
				break
			}
		}

		if characterToAssign == nil {
			http.Error(w, "Character not found in unassigned characters", http.StatusNotFound)
			return
		}

		// Assign to existing account or create a new one
		if request.AccountID != nil {
			// Assign to existing account
			var accountFound bool
			for i, account := range accounts {
				if account.ID == *request.AccountID {
					accounts[i].Characters = append(accounts[i].Characters, *characterToAssign)
					accountFound = true
					break
				}
			}
			if !accountFound {
				http.Error(w, "Account not found", http.StatusNotFound)
				return
			}
		} else {
			// Create a new account
			newAccount := model.Account{
				Name:       "New Account",
				Status:     "Alpha", // Or determine status appropriately
				Characters: []model.CharacterIdentity{*characterToAssign},
				ID:         time.Now().Unix(), // Generate a unique ID
			}
			accounts = append(accounts, newAccount)
		}

		// Remove the character from unassigned characters
		var updatedUnassigned []model.CharacterIdentity
		for _, character := range unassignedCharacters {
			if character.Character.CharacterID != request.CharacterID {
				updatedUnassigned = append(updatedUnassigned, character)
			}
		}

		// Save the updated accounts and unassigned characters
		err = persist.SaveAccounts(mainIdentity, accounts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		err = persist.SaveUnassignedCharacters(mainIdentity, updatedUnassigned)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func RemoveCharacterHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			CharacterID int64 `json:"characterID"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if request.CharacterID == 0 {
			http.Error(w, "CharacterID is required", http.StatusBadRequest)
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

		// Fetch accounts and unassigned characters
		accounts, err := persist.FetchAccountByIdentity(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		// Find and remove the character from accounts
		var characterToRemove *model.CharacterIdentity
		var accountIndex, charIndex int
		var characterFound bool

		for i := range accounts {
			for j := range accounts[i].Characters {
				if accounts[i].Characters[j].Character.CharacterID == request.CharacterID {
					characterToRemove = &accounts[i].Characters[j]
					accountIndex = i
					charIndex = j
					characterFound = true
					break
				}
			}
			if characterFound {
				break
			}
		}

		if !characterFound {
			http.Error(w, "Character not found in accounts", http.StatusNotFound)
			return
		}

		// Remove character from account
		accounts[accountIndex].Characters = append(
			accounts[accountIndex].Characters[:charIndex],
			accounts[accountIndex].Characters[charIndex+1:]...,
		)

		// Add character to unassigned characters
		unassignedCharacters = append(unassignedCharacters, *characterToRemove)

		// Save the updated accounts and unassigned characters
		err = persist.SaveAccounts(mainIdentity, accounts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		err = persist.SaveUnassignedCharacters(mainIdentity, unassignedCharacters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
