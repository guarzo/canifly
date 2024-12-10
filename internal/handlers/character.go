package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
	"github.com/sirupsen/logrus"
)

// CharacterHandler manages character-related API endpoints.
type CharacterHandler struct {
	sessionService *flyHttp.SessionService
	logger         *logrus.Logger
	configService  *services.ConfigService
	datastore      *persist.DataStore
}

// NewCharacterHandler constructs a CharacterHandler with the given dependencies.
func NewCharacterHandler(s *flyHttp.SessionService, l *logrus.Logger, c *services.ConfigService, d *persist.DataStore) *CharacterHandler {
	return &CharacterHandler{
		sessionService: s,
		logger:         l,
		configService:  c,
		datastore:      d,
	}
}

// respondJSON sends a success response with JSON-encoded data.
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// respondError is a helper to send an error response in JSON format.
func respondError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, fmt.Sprintf(`{"error":"%s"}`, msg), code)
}

// decodeJSONBody decodes the JSON body into the provided dst.
func decodeJSONBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

// UpdateCharacter updates fields for a given character such as Role or MCT flag.
func (h *CharacterHandler) UpdateCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CharacterID int64                  `json:"characterID"`
			Updates     map[string]interface{} `json:"updates"`
		}

		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.CharacterID == 0 || len(request.Updates) == 0 {
			respondError(w, "CharacterID and updates are required", http.StatusBadRequest)
			return
		}

		accounts, err := h.datastore.FetchAccounts()
		if err != nil {
			respondError(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		// Find the character to update
		charIdentity := h.findCharacterInAccounts(accounts, request.CharacterID)
		if charIdentity == nil {
			respondError(w, "Character not found", http.StatusNotFound)
			return
		}

		// Apply updates
		for key, value := range request.Updates {
			switch key {
			case "Role":
				roleStr, ok := value.(string)
				if !ok {
					respondError(w, "Role must be a string", http.StatusBadRequest)
					return
				}
				if err := h.configService.UpdateRoles(roleStr); err != nil {
					h.logger.Infof("Failed to update roles: %v", err)
				}
				charIdentity.Role = roleStr

			case "MCT":
				mctBool, ok := value.(bool)
				if !ok {
					respondError(w, "MCT must be boolean", http.StatusBadRequest)
					return
				}
				charIdentity.MCT = mctBool

			default:
				respondError(w, fmt.Sprintf("Unknown update field: %s", key), http.StatusBadRequest)
				return
			}
		}

		if err := h.datastore.SaveAccounts(accounts); err != nil {
			respondError(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// RemoveCharacter removes a character from the accounts.
func (h *CharacterHandler) RemoveCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CharacterID int64 `json:"characterID"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if request.CharacterID == 0 {
			respondError(w, "CharacterID is required", http.StatusBadRequest)
			return
		}
		h.logger.Infof("Beginning character removal for %d", request.CharacterID)

		accounts, err := h.datastore.FetchAccounts()
		if err != nil {
			respondError(w, fmt.Sprintf("Error fetching accounts: %v", err), http.StatusInternalServerError)
			return
		}

		accountIndex, charIndex, found := h.findCharacterIndices(accounts, request.CharacterID)
		if !found {
			respondError(w, "Character not found in accounts", http.StatusNotFound)
			return
		}
		h.logger.Infof("before removal length of %s characters is %d", accounts[accountIndex].Name, len(accounts[accountIndex].Characters))

		// Remove the character from the account
		accounts[accountIndex].Characters = append(
			accounts[accountIndex].Characters[:charIndex],
			accounts[accountIndex].Characters[charIndex+1:]...,
		)

		accountName := accounts[accountIndex].Name

		h.logger.Infof("after removal - length of %s characters is %d", accountName, len(accounts[accountIndex].Characters))

		if err := h.datastore.SaveAccounts(accounts); err != nil {
			respondError(w, fmt.Sprintf("Failed to save accounts: %v", err), http.StatusInternalServerError)
			return
		}

		accounts, err = h.datastore.FetchAccounts()
		if err != nil {
			h.logger.Infof("unable to retrieve accounts for verification")
		} else {
			index, found := findAccountIndex(accounts, accountName)
			if found {
				h.logger.Infof("after removal - length of %s characters is %d", accountName, len(accounts[index].Characters))
			}
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// AssociateCharacter associates a character with a user ID and updates config.
func (h *CharacterHandler) AssociateCharacter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   string `json:"userId"`
		CharId   string `json:"charId"`
		UserName string `json:"userName"`
		CharName string `json:"charName"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.configService.AssociateCharacter(req.UserId, req.CharId); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%s associated with %s", req.CharName, req.UserName),
	})
}

// UnassociateCharacter removes an association between a user and a character.
func (h *CharacterHandler) UnassociateCharacter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   string `json:"userId"`
		CharId   string `json:"charId"`
		UserName string `json:"userName"`
		CharName string `json:"charName"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.configService.UnassociateCharacter(req.UserId, req.CharId); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%s has been unassociated from %s", req.CharName, req.UserName),
	})
}

// SaveUserSelections saves user selections for directories/profiles.
func (h *CharacterHandler) SaveUserSelections(w http.ResponseWriter, r *http.Request) {
	var req model.UserSelections
	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.configService.UpdateUserSelections(req); err != nil {
		respondError(w, fmt.Sprintf("Failed to save user selections: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}

// SyncSubDirectory performs syncing of user/char files for a specific subdirectory.
func (h *CharacterHandler) SyncSubDirectory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userFilesCopied, charFilesCopied, err := h.configService.SyncDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync %v", err)})
		return
	}

	message := fmt.Sprintf("Synchronization complete in \"%s\", %d user files and %d character files copied.", req.SubDir, userFilesCopied, charFilesCopied)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// SyncAllSubdirectories syncs user/char files across all subdirectories.
func (h *CharacterHandler) SyncAllSubdirectories(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userFilesCopied, charFilesCopied, err := h.configService.SyncAllDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync all: %v", err)})
		return
	}

	message := fmt.Sprintf("Sync completed for all subdirectories: %d user files copied and %d character files copied, based on user/char files from \"%s\".",
		userFilesCopied, charFilesCopied, req.SubDir)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// BackupDirectory creates a backup of the current settings directory.
func (h *CharacterHandler) BackupDirectory(w http.ResponseWriter, r *http.Request) {
	success, message := h.configService.BackupDir()
	respondJSON(w, map[string]interface{}{"success": success, "message": message})
}

// ChooseSettingsDir updates the settings directory to a user-specified path.
func (h *CharacterHandler) ChooseSettingsDir(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Directory string `json:"directory"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Directory == "" {
		respondError(w, "Directory is required", http.StatusBadRequest)
		return
	}

	if err := h.configService.UpdateSettingsDir(req.Directory); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{"success": true, "settingsDir": req.Directory})
}

// findCharacterInAccounts returns a pointer to the character identity if found.
func (h *CharacterHandler) findCharacterInAccounts(accounts []model.Account, characterID int64) *model.CharacterIdentity {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				return &accounts[i].Characters[j]
			}
		}
	}
	return nil
}

// findCharacterIndices returns indices of the account and character if found.
func (h *CharacterHandler) findCharacterIndices(accounts []model.Account, characterID int64) (int, int, bool) {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}
