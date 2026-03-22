// Package pagination provides generic GORM-based pagination with centered page navigation.
package pagination

import (
	"math"

	"gorm.io/gorm"
)

const (
	// DefaultLimit is the default number of items per page.
	DefaultLimit = 50
	// MaxLimit is the maximum allowed items per page.
	MaxLimit = 100
)

// Pagination holds the state for a paginated query result.
type Pagination[T any] struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`

	// Template-only fields (HTML pagination)
	Pages         []int `json:"-"` // centered pages
	ShowStartDots bool  `json:"-"`
	ShowEndDots   bool  `json:"-"`

	Items []T `json:"items"`
}

// ---------- Helpers ----------

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

// GetPage returns the current page number, defaulting to 1.
func (p *Pagination[T]) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

// GetOffset returns the offset for the current page.
func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.getLimit()
}

// ---------- Centered Pagination Logic ----------

func (p *Pagination[T]) buildCenteredPages(window int) {
	if p.TotalPages <= 1 {
		return
	}

	start := p.Page - window
	end := p.Page + window

	if start < 2 {
		start = 2
	}
	if end > p.TotalPages-1 {
		end = p.TotalPages - 1
	}

	for i := start; i <= end; i++ {
		p.Pages = append(p.Pages, i)
	}

	p.ShowStartDots = start > 2
	p.ShowEndDots = end < p.TotalPages-1
}

// ---------- GORM Scope ----------

// Paginate returns a GORM scope that applies pagination to a query.
func Paginate[T any](model any, p *Pagination[T]) func(db *gorm.DB) *gorm.DB {
	p.normalize()

	return func(db *gorm.DB) *gorm.DB {
		var total int64

		db.Session(&gorm.Session{}).
			Model(model).
			Count(&total)

		p.TotalItems = total
		p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))

		p.HasPrev = p.Page > 1
		p.HasNext = p.Page < p.TotalPages

		// 👇 centered window (recommended = 2)
		p.buildCenteredPages(2)

		return db.Offset(p.offSet()).Limit(p.getLimit())
	}
}
