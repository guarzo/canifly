package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type SettingsHandler struct {
	logger          interfaces.Logger
	settingsService interfaces.SettingsService
}

func NewSettingsHandler(
	l interfaces.Logger,
	setSvc interfaces.SettingsService,
) *SettingsHandler {
	return &SettingsHandler{
		logger:          l,
		settingsService: setSvc,
	}
}

// SaveUserSelections
func (h *SettingsHandler) SaveUserSelections(w http.ResponseWriter, r *http.Request) {
	var req model.UserSelections
	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.settingsService.SaveUserSelections(req); err != nil {
		respondError(w, fmt.Sprintf("Failed to save user selections: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}

// SyncSubDirectory
func (h *SettingsHandler) SyncSubDirectory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userFilesCopied, charFilesCopied, err := h.settingsService.SyncDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync %v", err)})
		return
	}

	message := fmt.Sprintf("Synchronization complete in \"%s\", %d user files and %d character files copied.",
		req.SubDir, userFilesCopied, charFilesCopied)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// SyncAllSubdirectories
func (h *SettingsHandler) SyncAllSubdirectories(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userFilesCopied, charFilesCopied, err := h.settingsService.SyncAllDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync all: %v", err)})
		return
	}

	message := fmt.Sprintf("Sync completed for all subdirectories: %d user files and %d character files copied, based on user/char files from \"%s\".",
		userFilesCopied, charFilesCopied, req.SubDir)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// BackupDirectory
func (h *SettingsHandler) BackupDirectory(w http.ResponseWriter, r *http.Request) {
	success, message := h.settingsService.BackupDir()
	respondJSON(w, map[string]interface{}{"success": success, "message": message})
}

// ChooseSettingsDir
func (h *SettingsHandler) ChooseSettingsDir(w http.ResponseWriter, r *http.Request) {
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

	if err := h.settingsService.UpdateSettingsDir(req.Directory); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{"success": true, "settingsDir": req.Directory})
}
