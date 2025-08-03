package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/guarzo/canifly/internal/services/fuzzworks"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type FuzzworksHandler struct {
	logger   interfaces.Logger
	basePath string
	cache    interfaces.HTTPCacheService
	wsHub    *WebSocketHub
}

func NewFuzzworksHandler(logger interfaces.Logger, basePath string, cache interfaces.HTTPCacheService, wsHub *WebSocketHub) *FuzzworksHandler {
	return &FuzzworksHandler{
		logger:   logger,
		basePath: basePath,
		cache:    cache,
		wsHub:    wsHub,
	}
}

// UpdateData handles POST /api/fuzzworks/update
func (h *FuzzworksHandler) UpdateData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Manual Fuzzworks data update requested")

		// Create a new Fuzzworks service instance for the update
		service := fuzzworks.New(h.logger, h.basePath, true) // force update

		// Set a timeout for the update operation
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		// Run the update
		updateStart := time.Now()
		err := service.Initialize(ctx)
		updateDuration := time.Since(updateStart)

		if err != nil {
			h.logger.Errorf("Fuzzworks update failed: %v", err)
			respondError(w, "Failed to update Fuzzworks data", http.StatusInternalServerError)
			return
		}

		// Invalidate any cached EVE data since we have new data
		InvalidateCache(h.cache, "eve:")

		// Broadcast update via WebSocket
		if h.wsHub != nil {
			h.wsHub.BroadcastUpdate("fuzzworks:updated", map[string]interface{}{
				"duration": updateDuration.Seconds(),
			})
		}

		h.logger.Infof("Fuzzworks data updated successfully in %.2f seconds", updateDuration.Seconds())

		respondJSON(w, map[string]interface{}{
			"success":  true,
			"duration": updateDuration.Seconds(),
			"message":  "Fuzzworks data updated successfully",
		})
	}
}

// GetStatus handles GET /api/fuzzworks/status
func (h *FuzzworksHandler) GetStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the metadata file path
		metadataPath := h.basePath + "/config/fuzzworks/fuzzworks_metadata.json"

		// Try to read the metadata
		var metadata struct {
			InvTypes     *fuzzworks.FileMetadata `json:"inv_types"`
			SolarSystems *fuzzworks.FileMetadata `json:"solar_systems"`
		}

		// Read metadata file if it exists
		if data, err := os.ReadFile(metadataPath); err == nil {
			json.Unmarshal(data, &metadata)
		}

		status := map[string]interface{}{
			"hasData": metadata.InvTypes != nil || metadata.SolarSystems != nil,
		}

		if metadata.InvTypes != nil {
			status["invTypes"] = map[string]interface{}{
				"lastUpdated": metadata.InvTypes.DownloadTime,
				"fileSize":    metadata.InvTypes.FileSize,
				"etag":        metadata.InvTypes.ETag,
			}
		}

		if metadata.SolarSystems != nil {
			status["solarSystems"] = map[string]interface{}{
				"lastUpdated": metadata.SolarSystems.DownloadTime,
				"fileSize":    metadata.SolarSystems.FileSize,
				"etag":        metadata.SolarSystems.ETag,
			}
		}

		respondJSON(w, status)
	}
}
