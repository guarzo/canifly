package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/persist"
)

func UpdateCharacterRoleHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			CharacterID int64  `json:"characterID"`
			Role        string `json:"role"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if request.CharacterID == 0 || request.Role == "" {
			http.Error(w, "CharacterID and Role are required", http.StatusBadRequest)
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

		// Fetch ConfigData
		configData, err := persist.FetchConfigData(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching config data: %v", err), http.StatusInternalServerError)
			return
		}

		// Update the roles list if new role
		roleExists := false
		for _, role := range configData.Roles {
			if role == request.Role {
				roleExists = true
				break
			}
		}
		if !roleExists {
			configData.Roles = append(configData.Roles, request.Role)
		}

		// Save updated ConfigData
		err = persist.SaveConfigData(mainIdentity, configData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save config data: %v", err), http.StatusInternalServerError)
			return
		}

		// Update the character's role in memory (since we're not persisting HomeData)
		// Fetch accounts (non-persisted)
		accounts, err := persist.FetchAccountByIdentity(mainIdentity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the character and update its role
		characterFound := false
		for i := range accounts {
			for j := range accounts[i].Characters {
				character := &accounts[i].Characters[j]
				if character.Character.CharacterID == request.CharacterID {
					character.Role = request.Role
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

		// Since accounts are not persisted here, no need to save them

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
