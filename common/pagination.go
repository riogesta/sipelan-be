package common

import (
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Data       any        `json:"data,omitempty"`
	Pagination Pagination `json:"pagination"`
}

func NewPaginatedResponse(message string, data any, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Status:     200,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

// GetPaginationParams extracts page and limit from query parameters.
// Defaults: page=1, limit=10. Max limit=100.
func GetPaginationParams(r *http.Request) (page int, limit int) {
	page = 1
	limit = 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100
			}
		}
	}

	return page, limit
}

// CalculateOffset returns the database offset for the given page and limit.
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages returns the total number of pages.
func CalculateTotalPages(totalItems int64, limit int) int {
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}
