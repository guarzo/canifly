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
