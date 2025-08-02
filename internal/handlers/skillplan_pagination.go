package handlers

import (
	"github.com/guarzo/canifly/internal/model"
)

// SkillPlanEntry represents a skill plan for pagination
type SkillPlanEntry struct {
	Name   string                    `json:"name"`
	Status model.SkillPlanWithStatus `json:"status"`
}

// PaginateSkillPlans applies pagination to skill plan map
func PaginateSkillPlans(skillPlans map[string]model.SkillPlanWithStatus, params PaginationParams) PaginatedResponse {
	// Convert map to slice for pagination
	entries := make([]SkillPlanEntry, 0, len(skillPlans))
	for name, status := range skillPlans {
		entries = append(entries, SkillPlanEntry{
			Name:   name,
			Status: status,
		})
	}
	
	total := len(entries)
	start := params.Offset
	end := start + params.Limit
	
	// Adjust bounds
	if start >= total {
		start = total
		end = total
	} else if end > total {
		end = total
	}
	
	// Get paginated slice
	var paginatedEntries []SkillPlanEntry
	if start < total {
		paginatedEntries = entries[start:end]
	} else {
		paginatedEntries = []SkillPlanEntry{}
	}
	
	// Calculate pagination metadata
	hasNext := end < total
	hasPrev := params.Page > 1
	
	return PaginatedResponse{
		Data: paginatedEntries,
		Pagination: PaginationParams{
			Page:    params.Page,
			Limit:   params.Limit,
			Offset:  params.Offset,
			HasNext: hasNext,
			HasPrev: hasPrev,
			Total:   total,
		},
	}
}