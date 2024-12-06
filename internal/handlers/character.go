package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"

	"github.com/sirupsen/logrus"
)

type CharacterHandler struct {
	sessionService *http2.SessionService
	logger         *logrus.Logger
	configService  *services.ConfigService
}

func NewCharacterHandler(s *http2.SessionService, l *logrus.Logger, c *services.ConfigService) *CharacterHandler {
	return &CharacterHandler{
		sessionService: s,
		logger:         l,
		configService:  c,
	}
}

func (h *CharacterHandler) UpdateCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CharacterID int64                  `json:"characterID"`
			Updates     map[string]interface{} `json:"updates"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.CharacterID == 0 || len(request.Updates) == 0 {
			http.Error(w, "CharacterID and updates are required", http.StatusBadRequest)
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

		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		var characterFound bool
		for i := range accounts {
			for j := range accounts[i].Characters {
				charIdentity := &accounts[i].Characters[j]
				if charIdentity.Character.CharacterID == request.CharacterID {
					for key, value := range request.Updates {
						switch key {
						case "Role":
							if roleStr, ok := value.(string); ok {
								err = h.configService.UpdateRoles(roleStr)
								if err != nil {
									h.logger.Info(err)
								}
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

		if err = persist.SaveAccounts(accounts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func (h *CharacterHandler) GetUnassignedCharacters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("start of get unassigned characters handler")

		session, _ := h.sessionService.Get(r, http2.SessionName)

		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
		if !ok || loggedIn == 0 {
			h.logger.Error("Attempt to get unassigned characters without main identity")
			handleErrorWithRedirect(w, r, "/logout")
			return
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters()
		if err != nil {
			http.Error(w, "Error fetching unassigned characters", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(unassignedCharacters); err != nil {
			http.Error(w, "Error encoding unassigned characters", http.StatusInternalServerError)
		}
	}
}

func (h *CharacterHandler) AssignCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CharacterID int64  `json:"characterID"`
			AccountID   *int64 `json:"accountID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
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

		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

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

		if request.AccountID != nil {
			var accountFound bool
			for i := range accounts {
				if accounts[i].ID == *request.AccountID {
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
			newAccount := model.Account{
				Name:       "New Account",
				Status:     "Alpha",
				Characters: []model.CharacterIdentity{*characterToAssign},
				ID:         time.Now().Unix(),
			}
			accounts = append(accounts, newAccount)
		}

		var updatedUnassigned []model.CharacterIdentity
		for _, character := range unassignedCharacters {
			if character.Character.CharacterID != request.CharacterID {
				updatedUnassigned = append(updatedUnassigned, character)
			}
		}

		if err = persist.SaveAccounts(accounts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}
		if err = persist.SaveUnassignedCharacters(updatedUnassigned); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

func (h *CharacterHandler) RemoveCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CharacterID int64 `json:"characterID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if request.CharacterID == 0 {
			http.Error(w, "CharacterID is required", http.StatusBadRequest)
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

		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

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

		accounts[accountIndex].Characters = append(
			accounts[accountIndex].Characters[:charIndex],
			accounts[accountIndex].Characters[charIndex+1:]...,
		)

		unassignedCharacters = append(unassignedCharacters, *characterToRemove)

		if err = persist.SaveAccounts(accounts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}
		if err = persist.SaveUnassignedCharacters(unassignedCharacters); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save unassigned characters: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
