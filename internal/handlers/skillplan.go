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
			http.Error(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		content, err := h.skillService.GetSkillPlanFile(planName)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Skill plan file not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to read eve plan file: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, writeErr := w.Write(content); writeErr != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", writeErr), http.StatusInternalServerError)
		}
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
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if requestData.PlanName == "" {
			h.logger.Error("planName parameter missing")
			http.Error(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		if h.skillService.CheckIfDuplicatePlan(requestData.PlanName) {
			h.logger.Errorf("duplicate plan name %s", requestData.PlanName)
			http.Error(w, fmt.Sprintf("%s is already used as a plan name", requestData.PlanName), http.StatusBadRequest)
			return
		}

		if err := h.skillService.ParseAndSaveSkillPlan(requestData.Contents, requestData.PlanName); err != nil {
			h.logger.Errorf("Failed to save eve plan: %v", err)
			http.Error(w, "Failed to save eve plan", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

func (h *SkillPlanHandler) DeleteSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		planName := r.URL.Query().Get("planName")
		if planName == "" {
			h.logger.Error("planName parameter missing")
			http.Error(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		if err := h.skillService.DeleteSkillPlan(planName); err != nil {
			h.logger.Errorf("Failed to delete eve plan: %v", err)
			http.Error(w, "Failed to delete eve plan", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}
