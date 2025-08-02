package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type ConfigHandler struct {
	logger        interfaces.Logger
	configService interfaces.ConfigurationService
	cache         interfaces.HTTPCacheService
}

func NewConfigHandler(
	l interfaces.Logger,
	s interfaces.ConfigurationService,
	cache interfaces.HTTPCacheService,
) *ConfigHandler {
	return &ConfigHandler{
		logger:        l,
		configService: s,
		cache:         cache,
	}
}

// RESTful endpoint: GET /api/config
func (h *ConfigHandler) GetConfig() http.HandlerFunc {
	return WithCache(
		h.cache,
		h.logger,
		"config:data",
		10*time.Minute, // Cache for 10 minutes
		func() (interface{}, error) {
			config, err := h.configService.FetchConfigData()
			if err != nil {
				return nil, err
			}
			
			configData := map[string]interface{}{
				"settingsDir": config.SettingsDir,
				"roles": config.Roles,
				"userSelections": config.DropDownSelections,
				"lastBackupDir": config.LastBackupDir,
			}
			
			// Return in the format the frontend expects
			return map[string]interface{}{
				"data": configData,
			}, nil
		},
	)
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
		
		// Invalidate config cache after successful update
		InvalidateCache(h.cache, "config:")
		
		respondJSON(w, map[string]bool{"success": true})
	}
}
