package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type SkillPlanHandler struct {
	logger         interfaces.Logger
	eveDataService interfaces.EVEDataService
	accountService interfaces.AccountManagementService
	cache          interfaces.HTTPCacheService
	wsHub          *WebSocketHub
}

func NewSkillPlanHandler(l interfaces.Logger, e interfaces.EVEDataService, a interfaces.AccountManagementService, cache interfaces.HTTPCacheService, wsHub *WebSocketHub) *SkillPlanHandler {
	return &SkillPlanHandler{
		logger:         l,
		eveDataService: e,
		accountService: a,
		cache:          cache,
		wsHub:          wsHub,
	}
}

// RESTful endpoints

// ListSkillPlans handles GET /api/skill-plans
func (h *SkillPlanHandler) ListSkillPlans() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse pagination parameters
		paginationParams := ParsePaginationParams(r)

		// Check cache first (cache key includes pagination params)
		cacheKey := fmt.Sprintf("skillplans:list:page:%d:limit:%d", paginationParams.Page, paginationParams.Limit)

		cacheHandler := WithCache(
			h.cache,
			h.logger,
			cacheKey,
			3*time.Minute, // Cache for 3 minutes
			func() (interface{}, error) {
				// Get accounts to calculate skill plan status
				accounts, err := h.accountService.FetchAccounts()
				if err != nil {
					return nil, err
				}

				// Get skill plans with calculated status
				skillPlans, _ := h.eveDataService.GetPlanAndConversionData(
					accounts,
					h.eveDataService.GetSkillPlans(),
					h.eveDataService.GetSkillTypes(),
				)

				// Apply pagination to skill plans
				paginatedResponse := PaginateSkillPlans(skillPlans, paginationParams)

				return paginatedResponse, nil
			},
		)
		cacheHandler(w, r)
	}
}

// GetSkillPlan handles GET /api/skill-plans/{name}
func (h *SkillPlanHandler) GetSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		planName := vars["name"]

		if planName == "" {
			respondError(w, "Missing plan name", http.StatusBadRequest)
			return
		}

		content, err := h.eveDataService.GetSkillPlanFile(planName)
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, fmt.Sprintf("Skill plan %s not found", planName), http.StatusNotFound)
			} else {
				h.logger.Errorf("Failed to get skill plan %s: %v", planName, err)
				respondError(w, "Failed to get skill plan", http.StatusInternalServerError)
			}
			return
		}

		respondJSON(w, map[string]interface{}{
			"name":    planName,
			"content": content,
		})
	}
}

// CreateSkillPlan handles POST /api/skill-plans
func (h *SkillPlanHandler) CreateSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Name == "" || request.Content == "" {
			respondError(w, "Name and content are required", http.StatusBadRequest)
			return
		}

		if h.eveDataService.CheckIfDuplicatePlan(request.Name) {
			respondError(w, fmt.Sprintf("Skill plan %s already exists", request.Name), http.StatusConflict)
			return
		}

		if err := h.eveDataService.ParseAndSaveSkillPlan(request.Content, request.Name); err != nil {
			h.logger.Errorf("Failed to create skill plan: %v", err)
			respondError(w, "Failed to create skill plan", http.StatusInternalServerError)
			return
		}

		// Invalidate skill plans cache after successful creation
		InvalidateCache(h.cache, "skillplans:")
		InvalidateCache(h.cache, "eve:skillplans")

		// Broadcast creation via WebSocket
		if h.wsHub != nil {
			h.wsHub.BroadcastUpdate("skillplan:created", map[string]interface{}{
				"name": request.Name,
			})
		}

		w.WriteHeader(http.StatusCreated)
		respondJSON(w, map[string]interface{}{
			"name":    request.Name,
			"content": request.Content,
		})
	}
}

// UpdateSkillPlan handles PUT /api/skill-plans/{name}
func (h *SkillPlanHandler) UpdateSkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		planName := vars["name"]

		var request struct {
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Content == "" {
			respondError(w, "Content is required", http.StatusBadRequest)
			return
		}

		// Check if plan exists
		_, err := h.eveDataService.GetSkillPlanFile(planName)
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, fmt.Sprintf("Skill plan %s not found", planName), http.StatusNotFound)
			} else {
				h.logger.Errorf("Failed to check skill plan %s: %v", planName, err)
				respondError(w, "Failed to update skill plan", http.StatusInternalServerError)
			}
			return
		}

		// Delete and recreate (since the service doesn't have an update method)
		if err := h.eveDataService.DeleteSkillPlan(planName); err != nil {
			h.logger.Errorf("Failed to delete old skill plan: %v", err)
			respondError(w, "Failed to update skill plan", http.StatusInternalServerError)
			return
		}

		if err := h.eveDataService.ParseAndSaveSkillPlan(request.Content, planName); err != nil {
			h.logger.Errorf("Failed to save updated skill plan: %v", err)
			respondError(w, "Failed to update skill plan", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"name":    planName,
			"content": request.Content,
		})
	}
}

// DeleteSkillPlanRESTful handles DELETE /api/skill-plans/{name}
func (h *SkillPlanHandler) DeleteSkillPlanRESTful() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		planName := vars["name"]

		if planName == "" {
			respondError(w, "Missing plan name", http.StatusBadRequest)
			return
		}

		if err := h.eveDataService.DeleteSkillPlan(planName); err != nil {
			if os.IsNotExist(err) {
				respondError(w, fmt.Sprintf("Skill plan %s not found", planName), http.StatusNotFound)
			} else {
				h.logger.Errorf("Failed to delete skill plan: %v", err)
				respondError(w, "Failed to delete skill plan", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// CopySkillPlan handles POST /api/skill-plans/{name}/copy
func (h *SkillPlanHandler) CopySkillPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sourceName := vars["name"]

		var request struct {
			NewName string `json:"newName"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.NewName == "" {
			respondError(w, "New name is required", http.StatusBadRequest)
			return
		}

		// Get source plan
		content, err := h.eveDataService.GetSkillPlanFile(sourceName)
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, fmt.Sprintf("Source skill plan %s not found", sourceName), http.StatusNotFound)
			} else {
				h.logger.Errorf("Failed to get source skill plan: %v", err)
				respondError(w, "Failed to copy skill plan", http.StatusInternalServerError)
			}
			return
		}

		// Check if new name already exists
		if h.eveDataService.CheckIfDuplicatePlan(request.NewName) {
			respondError(w, fmt.Sprintf("Skill plan %s already exists", request.NewName), http.StatusConflict)
			return
		}

		// Create copy
		if err := h.eveDataService.ParseAndSaveSkillPlan(string(content), request.NewName); err != nil {
			h.logger.Errorf("Failed to create skill plan copy: %v", err)
			respondError(w, "Failed to copy skill plan", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		respondJSON(w, map[string]interface{}{
			"name":    request.NewName,
			"content": content,
		})
	}
}

// RefreshSkillPlans handles POST /api/skill-plans/refresh
func (h *SkillPlanHandler) RefreshSkillPlans() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Refresh skill plans from remote repository
		if err := h.eveDataService.RefreshRemotePlans(); err != nil {
			h.logger.Errorf("Failed to refresh skill plans: %v", err)
			respondError(w, "Failed to refresh skill plans", http.StatusInternalServerError)
			return
		}

		// Clear cache for skill plans
		InvalidateCache(h.cache, "skillplans:")

		respondJSON(w, map[string]string{
			"status":  "success",
			"message": "Skill plans refreshed successfully",
		})
	}
}
