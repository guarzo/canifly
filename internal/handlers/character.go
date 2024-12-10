package handlers

import (
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type CharacterHandler struct {
	logger           interfaces.Logger
	characterService interfaces.CharacterService
}

func NewCharacterHandler(
	l interfaces.Logger,
	c interfaces.CharacterService,
) *CharacterHandler {
	return &CharacterHandler{
		logger:           l,
		characterService: c,
	}
}

func (h *CharacterHandler) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
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

	if err := h.characterService.UpdateCharacterFields(request.CharacterID, request.Updates); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}

// RemoveCharacter removes a character via characterService
func (h *CharacterHandler) RemoveCharacter(w http.ResponseWriter, r *http.Request) {
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

	if err := h.characterService.RemoveCharacter(request.CharacterID); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}
