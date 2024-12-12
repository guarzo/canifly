package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type ConfigHandler struct {
	logger        interfaces.Logger
	configService interfaces.ConfigService
}

func NewConfigHandler(
	l interfaces.Logger,
	s interfaces.ConfigService,
) *ConfigHandler {
	return &ConfigHandler{
		logger:        l,
		configService: s,
	}
}

// SaveUserSelections
func (h *ConfigHandler) SaveUserSelections(w http.ResponseWriter, r *http.Request) {
	var req model.DropDownSelections
	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.configService.SaveUserSelections(req); err != nil {
		respondError(w, fmt.Sprintf("Failed to save user selections: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}

func (h *ConfigHandler) ChooseSettingsDir(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Directory string `json:"directory"`
	}
	h.logger.Infof("in choose settings handler")

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

func (h *ConfigHandler) ResetToDefaultDir(w http.ResponseWriter, r *http.Request) {

	h.logger.Infof("in reset to default dir handler")

	if err := h.configService.EnsureSettingsDir(); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{"success": true})
}
