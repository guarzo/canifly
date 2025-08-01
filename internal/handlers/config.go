package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type ConfigHandler struct {
	logger        interfaces.Logger
	configService interfaces.ConfigurationService
}

func NewConfigHandler(
	l interfaces.Logger,
	s interfaces.ConfigurationService,
) *ConfigHandler {
	return &ConfigHandler{
		logger:        l,
		configService: s,
	}
}

// RESTful endpoint: GET /api/config
func (h *ConfigHandler) GetConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config, err := h.configService.FetchConfigData()
		if err != nil {
			respondError(w, "Failed to fetch config", http.StatusInternalServerError)
			return
		}
		
		response := map[string]interface{}{
			"settingsDir": config.SettingsDir,
			"roles": config.Roles,
			"userSelections": config.DropDownSelections,
			"lastBackupDir": config.LastBackupDir,
		}
		
		respondJSON(w, response)
	}
}

// RESTful endpoint: PATCH /api/config
func (h *ConfigHandler) UpdateConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			SettingsDir    *string                    `json:"settingsDir,omitempty"`
			UserSelections *model.DropDownSelections `json:"userSelections,omitempty"`
			Roles          *[]string                 `json:"roles,omitempty"`
		}
		
		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Update settings directory if provided
		if request.SettingsDir != nil {
			if err := h.configService.UpdateSettingsDir(*request.SettingsDir); err != nil {
				respondError(w, fmt.Sprintf("Failed to update settings directory: %v", err), http.StatusInternalServerError)
				return
			}
		}
		
		// Update user selections if provided
		if request.UserSelections != nil {
			if err := h.configService.SaveUserSelections(*request.UserSelections); err != nil {
				respondError(w, fmt.Sprintf("Failed to save user selections: %v", err), http.StatusInternalServerError)
				return
			}
		}
		
		// Update roles if provided
		if request.Roles != nil {
			if err := h.configService.SaveRoles(*request.Roles); err != nil {
				respondError(w, fmt.Sprintf("Failed to save roles: %v", err), http.StatusInternalServerError)
				return
			}
		}
		
		respondJSON(w, map[string]bool{"success": true})
	}
}

// Legacy endpoint (to be deprecated)
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

// Legacy endpoint (to be deprecated)
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

// Legacy endpoint (to be deprecated)
func (h *ConfigHandler) ResetToDefaultDir(w http.ResponseWriter, r *http.Request) {

	h.logger.Infof("in reset to default dir handler")

	if err := h.configService.EnsureSettingsDir(); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{"success": true})
}
