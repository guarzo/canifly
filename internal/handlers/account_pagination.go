package handlers

import (
	"github.com/guarzo/canifly/internal/model"
)

// PaginateAccounts applies pagination to account slice
func PaginateAccounts(accounts []model.Account, params PaginationParams) PaginatedResponse {
	total := len(accounts)
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
	var paginatedAccounts []model.Account
	if start < total {
		paginatedAccounts = accounts[start:end]
	} else {
		paginatedAccounts = []model.Account{}
	}
	
	// Calculate pagination metadata  
	hasNext := end < total
	hasPrev := params.Page > 1
	
	return PaginatedResponse{
		Data: paginatedAccounts,
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