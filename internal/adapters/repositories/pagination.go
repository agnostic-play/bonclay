package repositories

import (
	"berlin.allobank.com/tools/bonclay/internal/common/pagination"
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
		if q == nil {
			return db
		}

		query := db

		// (field1 LIKE ? OR field2 LIKE ? OR ...)
		if q.Search != "" && len(fields) > 0 {
			term := "%" + q.Search + "%"
			query = query.Where(func(tx *gorm.DB) *gorm.DB {
				sub := tx
				for i, f := range fields {
					cond := f + " LIKE ?"
					if i == 0 {
						sub = sub.Where(cond, term)
					} else {
						sub = sub.Or(cond, term)
					}
				}
				return sub
			})
		}

		// AND filters...
		if q.Filters != nil && len(q.Filters) > 0 {
			for field, val := range q.Filters {
				// field should include operator placeholders if needed, e.g. "status = ?" / "created_at >= ?"
				query = query.Where(field, val)
			}
		}

		return query
	}
}
