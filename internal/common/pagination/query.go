package pagination

import (
	"strings"
)

type ListQuery struct {
	// Pagination controls
	Page     int `json:"page" form:"page"`           // Current page number (1-indexed)
	PageSize int `json:"page_size" form:"page_size"` // Items per page
	Limit    int `json:"limit" form:"limit"`         // Alternative: direct limit
	Offset   int `json:"offset" form:"offset"`       // Alternative: direct offset

	// Search and filtering
	Search string `json:"search" form:"search"` // Search query

	// Sorting
	SortBy    string `json:"sort_by" form:"sort_by"`       // Field to sort by
	SortOrder string `json:"sort_order" form:"sort_order"` // "asc" or "desc"
}

// Normalize validates and normalizes the query parameters
func (q *ListQuery) Normalize() {
	// Default page size
	if q.PageSize <= 0 && q.Limit <= 0 {
		q.PageSize = 10
	}

	// Convert between pagination styles
	if q.Page > 0 && q.PageSize > 0 {
		// Page-based pagination
		if q.Page < 1 {
			q.Page = 1
		}
		q.Limit = q.PageSize
		q.Offset = (q.Page - 1) * q.PageSize
	} else if q.Limit > 0 {
		// Offset-based pagination
		if q.Offset < 0 {
			q.Offset = 0
		}
		q.PageSize = q.Limit
		q.Page = (q.Offset / q.Limit) + 1
	}

	// Normalize sort order
	if q.SortOrder != "" && q.SortOrder != "asc" && q.SortOrder != "desc" {
		q.SortOrder = "asc"
	}

	// Trim search string
	q.Search = strings.TrimSpace(q.Search)
}

// GetLimit returns the effective limit for the query
func (q *ListQuery) GetLimit() int {
	return q.Limit
}

// GetOffset returns the effective offset for the query
func (q *ListQuery) GetOffset() int {
	return q.Offset
}
