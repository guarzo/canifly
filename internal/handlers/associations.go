package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AssociationHandler struct {
	logger       interfaces.Logger
	assocService interfaces.AssociationService
}

func NewAssociationHandler(
	l interfaces.Logger,
	a interfaces.AssociationService,
) *AssociationHandler {
	return &AssociationHandler{
		logger:       l,
		assocService: a,
	}
}

func (h *AssociationHandler) AssociateCharacter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   string `json:"userId"`
		CharId   string `json:"charId"`
		UserName string `json:"userName"`
		CharName string `json:"charName"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Infof("%v", req)

	if err := h.assocService.AssociateCharacter(req.UserId, req.CharId); err != nil {
		h.logger.Errorf("%v", err)
		respondJSON(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%s associated with %s", req.CharName, req.UserName),
	})
}

func (h *AssociationHandler) UnassociateCharacter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   string `json:"userId"`
		CharId   string `json:"charId"`
		UserName string `json:"userName"`
		CharName string `json:"charName"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.assocService.UnassociateCharacter(req.UserId, req.CharId); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%s has been unassociated from %s", req.CharName, req.UserName),
	})

}
