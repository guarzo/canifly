package handlers

import (
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters from query string
type PaginationParams struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasNext bool `json:"hasNext"`
	HasPrev bool `json:"hasPrev"`
	Total   int  `json:"total"`
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationParams `json:"pagination"`
}

// ParsePaginationParams extracts pagination parameters from request
func ParsePaginationParams(r *http.Request) PaginationParams {
	page := 1
	limit := 20 // Default limit

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 10000 { // Max 10000 items per page
			limit = l
		}
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// PaginateSlice applies pagination to a slice and returns paginated response
// This generic version works with any slice type
func PaginateSlice[T any](data []T, params PaginationParams) PaginatedResponse {
	total := len(data)
	start := params.Offset
	end := start + params.Limit

	// Adjust bounds
	if start >= total {
		start = total
		end = total
	} else if end > total || params.Limit == 0 {
		end = total
	}

	// Get paginated slice
	var paginatedData interface{}
	if start < total {
		paginatedData = data[start:end]
	} else {
		paginatedData = []T{}
	}

	// Calculate pagination metadata
	hasNext := end < total
	hasPrev := params.Page > 1

	return PaginatedResponse{
		Data: paginatedData,
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
