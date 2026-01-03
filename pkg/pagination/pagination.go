package pagination

import (
	"math"

	"gorm.io/gorm"
)

const (
	DefaultLimit = 50
	MaxLimit     = 100
)

type Pagination[T any] struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
	Items      []T   `json:"items"`
}

// ---- Helpers ----
func (p *Pagination[T]) normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
}

func (p *Pagination[T]) offSet() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination[T]) getLimit() int {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.getLimit()
}

// Paginate ---- GORM Scope ----
// Paginate[T any] is a function which will give paginated response
func Paginate[T any](model any, p *Pagination[T]) func(db *gorm.DB) *gorm.DB {
	p.normalize()

	return func(db *gorm.DB) *gorm.DB {
		var total int64

		// clean COUNT query
		db.Session(&gorm.Session{}).Model(model).Count(&total)

		p.TotalItems = total
		p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))

		p.HasPrev = p.Page > 1
		p.HasNext = p.Page < p.TotalPages

		return db.Offset(p.offSet()).Limit(p.getLimit())
	}
}
