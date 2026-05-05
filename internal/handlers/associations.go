package handlers

import (
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
				"userId":   assoc.UserId,
				"charId":   assoc.CharId,
				"charName": assoc.CharName,
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
			CharName    string `json:"charName"`
		}

		if err := decodeJSONBody(r, &request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.UserID == "" || request.CharacterID == "" {
			respondError(w, "userId and characterId are required", http.StatusBadRequest)
			return
		}

		if err := h.assocService.AssociateCharacter(request.UserID, request.CharacterID, request.CharName); err != nil {
			h.logger.Errorf("Failed to create association: %v", err)
			respondError(w, "Failed to create association", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"success":  true,
			"message":  "Character associated.",
			"userId":   request.UserID,
			"charId":   request.CharacterID,
			"charName": request.CharName,
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

		respondJSON(w, map[string]interface{}{
			"success": true,
			"message": "Character unassociated.",
		})
	}
}
