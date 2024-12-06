package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
)

type SkillPlanHandler struct {
	logger       *logrus.Logger
	skillService *services.SkillService
}

func NewSkillPlanHandler(l *logrus.Logger, s *services.SkillService) *SkillPlanHandler {
	return &SkillPlanHandler{logger: l,
		skillService: s}
}

func (h *SkillPlanHandler) GetSkillPlanFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		planName := r.URL.Query().Get("planName")
		if planName == "" {
			http.Error(w, "Missing planName parameter", http.StatusBadRequest)
			return
		}

		planName += ".txt"
		h.logger.Infof("Attempting to serve skill plan file: %s", planName)

		skillPlanDir, err := persist.GetWritablePlansPath()
		if err != nil {
			h.logger.Errorf("Failed to retrieve skill plan directory: %v", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve skill plan directory: %v", err), http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join(skillPlanDir, planName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Skill plan file not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to read skill plan file: %v", err), http.StatusInternalServerError)
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

		skills := h.skillService.ParseSkillPlanContents(requestData.Contents)
		if err := persist.SaveSkillPlan(requestData.PlanName, skills); err != nil {
			h.logger.Errorf("Failed to save skill plan: %v", err)
			http.Error(w, "Failed to save skill plan", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
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

		if err := persist.DeleteSkillPlan(planName); err != nil {
			h.logger.Errorf("Failed to delete skill plan: %v", err)
			http.Error(w, "Failed to delete skill plan", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
