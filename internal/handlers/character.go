package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type CharacterHandler struct {
	logger           interfaces.Logger
	characterService interfaces.CharacterService
	esiAPIService    interfaces.ESIAPIService
	cache            interfaces.HTTPCacheService
}

func NewCharacterHandler(
	l interfaces.Logger,
	cs interfaces.CharacterService,
	esi interfaces.ESIAPIService,
	c interfaces.HTTPCacheService,
) *CharacterHandler {
	return &CharacterHandler{
		logger:           l,
		characterService: cs,
		esiAPIService:    esi,
		cache:            c,
	}
}

// RESTful endpoint: GET /api/characters/:id
func (h *CharacterHandler) GetCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterID := vars["id"]

		character, err := h.esiAPIService.GetCharacter(characterID)
		if err != nil {
			respondError(w, "Character not found", http.StatusNotFound)
			return
		}

		respondJSON(w, character)
	}
}

// RESTful endpoint: PATCH /api/characters/:id
func (h *CharacterHandler) UpdateCharacterRESTful() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		var updates map[string]interface{}
		if err := decodeJSONBody(r, &updates); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(updates) == 0 {
			respondError(w, "No updates provided", http.StatusBadRequest)
			return
		}

		if err := h.characterService.UpdateCharacterFields(characterID, updates); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

// RESTful endpoint: DELETE /api/characters/:id
func (h *CharacterHandler) DeleteCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		if err := h.characterService.RemoveCharacter(characterID); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// RESTful endpoint: POST /api/characters/:id/refresh
func (h *CharacterHandler) RefreshCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		// Use EVE data service to refresh the character
		updated, err := h.characterService.RefreshCharacterData(characterID)
		if err != nil {
			if err.Error() == "character not found" {
				respondError(w, "Character not found", http.StatusNotFound)
			} else if err.Error() == "token expired" {
				respondError(w, "Authentication token expired", http.StatusUnauthorized)
			} else {
				respondError(w, fmt.Sprintf("Failed to refresh character: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Clear account cache to ensure fresh data on next fetch
		if h.cache != nil {
			// Invalidate all account cache entries
			h.cache.Invalidate("accounts:list:*")
			h.logger.Info("Cleared account cache after character refresh")
		}

		respondJSON(w, map[string]interface{}{
			"success": true,
			"message": "Character refreshed successfully",
			"updated": updated,
		})
	}
}
