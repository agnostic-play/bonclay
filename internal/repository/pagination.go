package repository

import (
	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"gorm.io/gorm"
)

// Paginate returns a GORM scope for pagination
func Paginate(q *pagination.ListQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q.Normalize()
		return db.Offset(q.Offset).Limit(q.Limit)
	}
}

// Sort returns a GORM scope for sorting
func Sort(q *pagination.ListQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.SortBy != "" {
			order := q.SortBy
			if q.SortOrder != "" {
				order += " " + q.SortOrder
			}
			return db.Order(order)
		}
		return db
	}
}

// SearchScope returns a GORM scope for searching
// You need to provide the fields to search in
func SearchScope(q *pagination.ListQuery, fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.Search == "" || len(fields) == 0 {
			return db
		}

		query := db
		for i, field := range fields {
			searchTerm := "%" + q.Search + "%"
			if i == 0 {
				query = query.Where(field+" LIKE ?", searchTerm)
			} else {
				query = query.Or(field+" LIKE ?", searchTerm)
			}
		}

		if q.Filters != nil && len(q.Filters) > 0 {
			for field, val := range q.Filters {
				query = query.Where(field, val)
			}
		}

		return query
	}
}
