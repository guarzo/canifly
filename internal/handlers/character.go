package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type CharacterHandler struct {
	logger           interfaces.Logger
	eveDataService interfaces.EVEDataService
}

func NewCharacterHandler(
	l interfaces.Logger,
	e interfaces.EVEDataService,
) *CharacterHandler {
	return &CharacterHandler{
		logger:           l,
		eveDataService: e,
	}
}

// RESTful endpoint: GET /api/characters/:id
func (h *CharacterHandler) GetCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterID := vars["id"]

		character, err := h.eveDataService.GetCharacter(characterID)
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

		if err := h.eveDataService.UpdateCharacterFields(characterID, updates); err != nil {
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

		if err := h.eveDataService.RemoveCharacter(characterID); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Legacy endpoint (to be deprecated)
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

	if err := h.eveDataService.UpdateCharacterFields(request.CharacterID, request.Updates); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}

// Legacy endpoint (to be deprecated)
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

	if err := h.eveDataService.RemoveCharacter(request.CharacterID); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"success": true})
}
