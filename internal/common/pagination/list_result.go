package pagination

type ListResult[T any] struct {
	List       []T  `json:"list"`        // The actual items
	Total      int  `json:"total"`       // Total number of items
	Page       int  `json:"page"`        // Current page number
	PageSize   int  `json:"page_size"`   // Items per page
	TotalPages int  `json:"total_pages"` // Total number of pages
	HasNext    bool `json:"has_next"`    // Whether there's a next page
	HasPrev    bool `json:"has_prev"`    // Whether there's a previous page

	// Optional metadata
	Search    string `json:"search,omitempty"`     // Applied search query
	SortBy    string `json:"sort_by,omitempty"`    // Applied sort field
	SortOrder string `json:"sort_order,omitempty"` // Applied sort order
}

// NewListResult creates a new paginated result
func NewListResult[T any](query *ListQuery) *ListResult[T] {
	query.Normalize()
	totalPages := 0

	return &ListResult[T]{
		List:       []T{},
		Total:      0,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: 0,
		HasNext:    query.Page < totalPages,
		HasPrev:    query.Page > 1,
		Search:     query.Search,
		SortBy:     query.SortBy,
		SortOrder:  query.SortOrder,
	}
}

func (l *ListResult[T]) SetData(totalData int, data []T) {
	if l.PageSize > 0 {
		l.TotalPages = (totalData + l.PageSize - 1) / l.PageSize
	}

	l.Total = totalData
	l.List = data
}
