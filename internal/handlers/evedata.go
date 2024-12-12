package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type EveDataHandler struct {
	logger interfaces.Logger
	eveSvc interfaces.EveProfilesService
}

func NewEveDataHandler(
	l interfaces.Logger,
	s interfaces.EveProfilesService,
) *EveDataHandler {
	return &EveDataHandler{
		logger: l,
		eveSvc: s,
	}
}

// SyncSubDirectory
func (h *EveDataHandler) SyncSubDirectory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userFilesCopied, charFilesCopied, err := h.eveSvc.SyncDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync %v", err)})
		return
	}

	message := fmt.Sprintf("Synchronization complete in \"%s\", %d user files and %d character files copied.",
		req.SubDir, userFilesCopied, charFilesCopied)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// SyncAllSubdirectories
func (h *EveDataHandler) SyncAllSubdirectories(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubDir string `json:"subDir"`
		UserId string `json:"userId"`
		CharId string `json:"charId"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		h.logger.Errorf("Invalid request body for SyncAllSubdirectories: %v", err)
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Infof("SyncAllSubdirectories request: Profile=%s, UserId=%s, CharId=%s", req.SubDir, req.UserId, req.CharId)
	userFilesCopied, charFilesCopied, err := h.eveSvc.SyncAllDir(req.SubDir, req.CharId, req.UserId)
	if err != nil {
		h.logger.Errorf("Failed to sync all subdirectories from base %s (UserId=%s, CharId=%s): %v", req.SubDir, req.UserId, req.CharId, err)
		respondJSON(w, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to sync all: %v", err)})
		return
	}

	message := fmt.Sprintf("Sync completed for all subdirectories: %d user files and %d character files copied, based on user/char files from \"%s\".",
		userFilesCopied, charFilesCopied, req.SubDir)
	h.logger.Infof("SyncAllSubdirectories completed successfully for base %s (UserId=%s, CharId=%s): %s", req.SubDir, req.UserId, req.CharId, message)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

func (h *EveDataHandler) BackupDirectory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TargetDir string `json:"targetDir"`
		BackupDir string `json:"backupDir"`
	}
	if err := decodeJSONBody(r, &req); err != nil {
		h.logger.Errorf("Invalid request body for BackupDirectory: %v", err)
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Infof("Received backup request. TargetDir=%s, BackupDir=%s", req.TargetDir, req.BackupDir)

	if err := h.eveSvc.BackupDir(req.TargetDir, req.BackupDir); err != nil {
		h.logger.Errorf("Failed to backup settings from %s to %s: %v", req.TargetDir, req.BackupDir, err)
		respondError(w, fmt.Sprintf("Failed to backup settings: %v", err), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Backed up settings to %s", req.BackupDir)
	h.logger.Infof("Backup request successful. %s", message)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}
