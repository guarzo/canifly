package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/model"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Surrogate-Control", "no-store")

	state := fmt.Sprintf("main-%d", time.Now().UnixNano())
	xlog.Logf("Login handler - getting auth url")
	url := api.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthCharacterHandler(w http.ResponseWriter, r *http.Request) {
	state := fmt.Sprintf("character-%d", time.Now().UnixNano())
	xlog.Logf(state)
	url := api.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xlog.Logf("callback request received")

		// Get the authorization code and state
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		// Exchange the code for a token
		token, err := api.ExchangeCode(code)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to exchange token for code: %s, state: %s, %v", code, state, err), "/")
			return
		}

		// Get the user info
		user, err := api.GetUserInfo(token)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to get user info: %v", err), "/")
			return
		}

		// Start or retrieve the session
		session, _ := s.Get(r, sessionName)

		// Handle main identity state
		if state[:4] == "main" {
			session.Values[loggedInUser] = user.CharacterID
		}

		// Add the character to the authenticated list if not already present
		if _, ok := session.Values[allAuthenticatedCharacters].([]int64); ok {
			if !slices.Contains(session.Values[allAuthenticatedCharacters].([]int64), user.CharacterID) {
				session.Values[allAuthenticatedCharacters] = append(session.Values[allAuthenticatedCharacters].([]int64), user.CharacterID)
			}
		} else {
			session.Values[allAuthenticatedCharacters] = []int64{user.CharacterID}
		}

		// Retrieve the main identity from the session
		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, fmt.Sprintf("main identity not found, current session: %v", session.Values), "/logout")
			return
		}

		// Retrieve the unassigned characters for the logged-in user
		unassignedCharacters, err := persist.FetchUnassignedCharacters(mainIdentity)
		if err != nil {
			http.Error(w, "Error fetching unassigned characters", http.StatusInternalServerError)
			return
		}

		// Check if the character is already assigned to an account
		var characterAssigned bool
		err = persist.UpdateAccounts(mainIdentity, func(account *model.Account) error {
			for i := range account.Characters {
				if account.Characters[i].Character.CharacterID == user.CharacterID {
					// If the character is found in an account, we update the token
					account.Characters[i].Token = *token
					characterAssigned = true
					break
				}
			}
			return nil
		})

		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to update account model: %v", err), "/")
			return
		}

		// If the character is not assigned to any account, add it to unassigned characters
		if !characterAssigned {
			newCharacter := model.CharacterIdentity{
				Token: *token,
				Character: model.Character{
					BaseCharacterResponse: model.BaseCharacterResponse{
						CharacterID:   user.CharacterID,
						CharacterName: user.CharacterName,
					},
				},
			}

			// Add to the unassigned characters list
			unassignedCharacters = append(unassignedCharacters, newCharacter)

			// Save the updated unassigned characters to the file
			err := persist.SaveUnassignedCharacters(mainIdentity, unassignedCharacters)
			if err != nil {
				http.Error(w, "Error saving unassigned characters", http.StatusInternalServerError)
				return
			}

		}

		// Save the session state
		if err := session.Save(r, w); err != nil {
			xlog.Logf("Error saving session: %v", err)
		}

		// Redirect to the React app
		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
	}
}

func LogoutHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Get(r, sessionName)
		if err != nil {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		clearSession(s, w, r)
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		// Respond with JSON instead of redirect
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func ResetIdentitiesHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, sessionName)
		mainIdentity, ok := session.Values[loggedInUser].(int64)

		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, "Attempt to reset identities without a main identity", "/logout")
			return
		}

		err := persist.DeleteAccount(mainIdentity)
		if err != nil {
			xlog.Logf("Failed to delete identity %d: %v", mainIdentity, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
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
		// Retrieve the session
		session, err := s.Get(r, sessionName)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Error retrieving session: %v", err), "/logout")
			return
		}

		// Parse the incoming request (CharacterID)
		var request struct {
			CharacterID int64 `json:"characterID"`
		}
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Invalid request body: %v", err), "/")
			return
		}

		// Retrieve the logged-in user (main identity)
		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, "Main identity not found in session", "/logout")
			return
		}

		// Retrieve accounts for the logged-in user
		accounts, err := persist.FetchAccountByIdentity(mainIdentity)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Error fetching accounts: %v", err), "/")
			return
		}

		// Retrieve unassigned characters from the file
		unassignedCharacters, err := persist.FetchUnassignedCharacters(mainIdentity)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Error fetching unassigned characters: %v", err), "/")
			return
		}

		// Find the character to assign
		var characterToAssign *model.CharacterIdentity
		for i, character := range unassignedCharacters {
			if character.Character.CharacterID == request.CharacterID {
				// Found the character to assign
				characterToAssign = &unassignedCharacters[i]
				break
			}
		}

		if characterToAssign == nil {
			handleErrorWithRedirect(w, r, "Character not found in unassigned characters", "/")
			return
		}

		// If no account exists, create a new account for the user
		if len(accounts) == 0 {
			log.Printf("No accounts found for user %d, creating a new account", mainIdentity)
			newAccount := model.Account{
				Name:       "New Account",
				Status:     "Alpha",
				Characters: []model.CharacterIdentity{*characterToAssign},
				ID:         time.Now().Unix(),
			}
			accounts = append(accounts, newAccount)
		} else {
			// Assign the character to the first account (or implement account selection logic)
			accounts[0].Characters = append(accounts[0].Characters, *characterToAssign)
		}

		// Update account and unassigned characters
		err = persist.UpdateAccounts(mainIdentity, func(account *model.Account) error {
			// Add the character to the account
			account.Characters = append(account.Characters, *characterToAssign)

			// Remove the character from unassigned characters
			for i, character := range unassignedCharacters {
				if character.Character.CharacterID == request.CharacterID {
					// Remove character from the slice
					unassignedCharacters = append(unassignedCharacters[:i], unassignedCharacters[i+1:]...)
					break
				}
			}

			// Save the updated unassigned characters back to the file
			err = persist.SaveUnassignedCharacters(mainIdentity, unassignedCharacters)
			if err != nil {
				return fmt.Errorf("failed to save unassigned characters: %v", err)
			}

			return nil
		})

		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to update account: %v", err), "/")
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
