package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AssociationHandler struct {
	logger       interfaces.Logger
	assocService interfaces.AccountManagementService
}

func NewAssociationHandler(
	l interfaces.Logger,
	a interfaces.AccountManagementService,
) *AssociationHandler {
	return &AssociationHandler{
		logger:       l,
		assocService: a,
	}
}

// RESTful endpoints

// ListAssociations handles GET /api/associations
func (h *AssociationHandler) ListAssociations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		associations, err := h.assocService.GetAssociations()
		if err != nil {
			h.logger.Errorf("Failed to get associations: %v", err)
			respondError(w, "Failed to get associations", http.StatusInternalServerError)
			return
		}
		
		// Convert associations to response format
		data := make([]map[string]interface{}, len(associations))
		for i, assoc := range associations {
			data[i] = map[string]interface{}{
				"userId":      assoc.UserId,
				"characterId": assoc.CharId,
			}
		}
		
		respondJSON(w, map[string]interface{}{
			"data": data,
		})
	}
}

// CreateAssociation handles POST /api/associations
func (h *AssociationHandler) CreateAssociation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			UserID      string `json:"userId"`
			CharacterID string `json:"characterId"`
		}

		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.UserID == "" || request.CharacterID == "" {
			respondError(w, "userId and characterId are required", http.StatusBadRequest)
			return
		}

		if err := h.assocService.AssociateCharacter(request.UserID, request.CharacterID); err != nil {
			h.logger.Errorf("Failed to create association: %v", err)
			respondError(w, fmt.Sprintf("Failed to create association: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		respondJSON(w, map[string]interface{}{
			"userId":      request.UserID,
			"characterId": request.CharacterID,
		})
	}
}

// DeleteAssociation handles DELETE /api/associations/{userId}/{characterId}
func (h *AssociationHandler) DeleteAssociation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userId"]
		characterID := vars["characterId"]

		if userID == "" || characterID == "" {
			respondError(w, "Missing userId or characterId", http.StatusBadRequest)
			return
		}

		if err := h.assocService.UnassociateCharacter(userID, characterID); err != nil {
			h.logger.Errorf("Failed to delete association: %v", err)
			respondError(w, "Failed to delete association", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
