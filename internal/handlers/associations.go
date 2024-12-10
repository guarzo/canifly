package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AssociationHandler struct {
	logger         interfaces.Logger
	accountService interfaces.AccountService
}

func NewAssociationHandler(
	l interfaces.Logger,
	a interfaces.AccountService,
) *AssociationHandler {
	return &AssociationHandler{
		logger:         l,
		accountService: a,
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

	if err := h.accountService.AssociateCharacter(req.UserId, req.CharId); err != nil {
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

	if err := h.accountService.UnassociateCharacter(req.UserId, req.CharId); err != nil {
		respondJSON(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%s has been unassociated from %s", req.CharName, req.UserName),
	})

}
