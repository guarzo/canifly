package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type SkillPlanHandler struct {
	logger       interfaces.Logger
	skillService interfaces.SkillService
}

func NewSkillPlanHandler(l interfaces.Logger, s interfaces.SkillService) *SkillPlanHandler {
	return &SkillPlanHandler{
		logger:       l,
		skillService: s,
	}
}

func (h *SkillPlanHandler) GetSkillPlanFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		planName := r.URL.Query().Get("planName")
		if planName == "" {
			respondError(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		content, err := h.skillService.GetSkillPlanFile(planName)
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, fmt.Sprintf("skill plan %s not found", planName), http.StatusNotFound)
			} else {
				respondError(w, fmt.Sprintf("Failed to read skill plan file %s: %v", planName, err), http.StatusInternalServerError)

			}
			return
		}

		respondEncodedData(w, content)
	}
}

func (h *SkillPlanHandler) SaveSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			PlanName string `json:"name"`
			Contents string `json:"contents"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			h.logger.Errorf("Failed to parse JSON body: %v", err)
			respondError(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if requestData.PlanName == "" {
			h.logger.Error("planName parameter missing")
			respondError(w, "Missing planName", http.StatusBadRequest)
			return
		}

		if h.skillService.CheckIfDuplicatePlan(requestData.PlanName) {
			h.logger.Errorf("duplicate plan name %s", requestData.PlanName)
			respondError(w, fmt.Sprintf("%s is already used as a plan name", requestData.PlanName), http.StatusBadRequest)
			return
		}

		if err := h.skillService.ParseAndSaveSkillPlan(requestData.Contents, requestData.PlanName); err != nil {
			h.logger.Errorf("Failed to save eve plan: %v", err)
			respondError(w, "Failed to save eve plan", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

func (h *SkillPlanHandler) DeleteSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		planName := r.URL.Query().Get("planName")
		if planName == "" {
			h.logger.Error("planName parameter missing")
			respondError(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		if err := h.skillService.DeleteSkillPlan(planName); err != nil {
			h.logger.Errorf("Failed to delete eve plan: %v", err)
			respondError(w, "Failed to delete eve plan", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}
