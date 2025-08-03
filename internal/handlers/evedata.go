package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type EveDataHandler struct {
	logger           interfaces.Logger
	syncService      interfaces.SyncService
	configService    interfaces.ConfigurationService
	skillPlanService interfaces.SkillPlanService
	profileService   interfaces.ProfileService
	accountService   interfaces.AccountManagementService
	cache            interfaces.HTTPCacheService
}

func NewEveDataHandler(
	l interfaces.Logger,
	s interfaces.SyncService,
	configService interfaces.ConfigurationService,
	skillPlanService interfaces.SkillPlanService,
	profileService interfaces.ProfileService,
	accountService interfaces.AccountManagementService,
	cache interfaces.HTTPCacheService,
) *EveDataHandler {
	return &EveDataHandler{
		logger:           l,
		syncService:      s,
		configService:    configService,
		skillPlanService: skillPlanService,
		profileService:   profileService,
		accountService:   accountService,
		cache:            cache,
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

	userFilesCopied, charFilesCopied, err := h.syncService.SyncDirectory(req.SubDir, req.CharId, req.UserId)
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
	userFilesCopied, charFilesCopied, err := h.syncService.SyncAllDirectories(req.SubDir, req.CharId, req.UserId)
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

	if err := h.syncService.BackupDirectory(req.TargetDir, req.BackupDir); err != nil {
		h.logger.Errorf("Failed to backup settings from %s to %s: %v", req.TargetDir, req.BackupDir, err)
		respondError(w, fmt.Sprintf("Failed to backup settings: %v", err), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Backed up settings to %s", req.BackupDir)
	h.logger.Infof("Backup request successful. %s", message)
	respondJSON(w, map[string]interface{}{"success": true, "message": message})
}

// GetSkillPlans returns all skill plans
func (h *EveDataHandler) GetSkillPlans(w http.ResponseWriter, r *http.Request) {
	cacheHandler := WithCache(
		h.cache,
		h.logger,
		"eve:skillplans",
		5*time.Minute, // Cache for 5 minutes
		func() (interface{}, error) {
			h.logger.Debug("Getting skill plans")

			// Get accounts to calculate skill plan status
			accounts, err := h.accountService.FetchAccounts()
			if err != nil {
				return nil, err
			}

			// Get skill plans with calculated status
			skillPlans, _ := h.skillPlanService.GetPlanAndConversionData(
				accounts,
				h.skillPlanService.GetSkillPlans(),
				h.skillPlanService.GetSkillTypes(),
			)

			return skillPlans, nil
		},
	)
	cacheHandler(w, r)
}

// GetEveProfiles returns all EVE profiles
func (h *EveDataHandler) GetEveProfiles(w http.ResponseWriter, r *http.Request) {
	cacheHandler := WithCache(
		h.cache,
		h.logger,
		"eve:profiles",
		10*time.Minute, // Cache for 10 minutes
		func() (interface{}, error) {
			h.logger.Debug("Getting EVE profiles")

			// Load character settings (EVE profiles)
			eveProfiles, err := h.profileService.LoadCharacterSettings()
			if err != nil {
				h.logger.Errorf("Failed to load character settings: %v", err)
				return nil, err // Return nil and error so WithCache can handle it properly
			}

			return eveProfiles, nil
		},
	)
	cacheHandler(w, r)
}

// GetEveConversions returns all EVE conversions
func (h *EveDataHandler) GetEveConversions(w http.ResponseWriter, r *http.Request) {
	cacheHandler := WithCache(
		h.cache,
		h.logger,
		"eve:conversions",
		5*time.Minute, // Cache for 5 minutes
		func() (interface{}, error) {
			h.logger.Debug("Getting EVE conversions")

			// Get accounts to calculate conversions
			accounts, err := h.accountService.FetchAccounts()
			if err != nil {
				return nil, err
			}

			// Get conversions
			_, eveConversions := h.skillPlanService.GetPlanAndConversionData(
				accounts,
				h.skillPlanService.GetSkillPlans(),
				h.skillPlanService.GetSkillTypes(),
			)

			return eveConversions, nil
		},
	)
	cacheHandler(w, r)
}
