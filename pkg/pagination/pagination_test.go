package pagination

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		wantPage  int
		wantLimit int
	}{
		{"zero values default", 0, 0, 1, DefaultLimit},
		{"negative values default", -1, -5, 1, DefaultLimit},
		{"limit capped at max", 1, 200, 1, MaxLimit},
		{"valid values unchanged", 3, 25, 3, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination[string]{Page: tt.page, Limit: tt.limit}
			p.normalize()
			if p.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", p.Page, tt.wantPage)
			}
			if p.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", p.Limit, tt.wantLimit)
			}
		})
	}
}

func TestGetPageAndGetOffset(t *testing.T) {
	p := &Pagination[string]{Page: 3, Limit: 10}
	if got := p.GetPage(); got != 3 {
		t.Errorf("GetPage() = %d, want 3", got)
	}
	if got := p.GetOffset(); got != 20 {
		t.Errorf("GetOffset() = %d, want 20", got)
	}
}

func TestGetPageDefaultsToOne(t *testing.T) {
	p := &Pagination[string]{Page: 0}
	if got := p.GetPage(); got != 1 {
		t.Errorf("GetPage() = %d, want 1", got)
	}
}

func TestBuildCenteredPages(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		totalPages    int
		window        int
		wantPages     []int
		wantStartDots bool
		wantEndDots   bool
	}{
		{"single page", 1, 1, 2, nil, false, false},
		{"page 1 of 10", 1, 10, 2, []int{2, 3}, false, true},
		{"page 5 of 10", 5, 10, 2, []int{3, 4, 5, 6, 7}, true, true},
		{"page 9 of 10", 9, 10, 2, []int{7, 8, 9}, true, false},
		{"page 2 of 3", 2, 3, 2, []int{2}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination[string]{Page: tt.page, TotalPages: tt.totalPages}
			p.buildCenteredPages(tt.window)

			if !intSliceEqual(p.Pages, tt.wantPages) {
				t.Errorf("Pages = %v, want %v", p.Pages, tt.wantPages)
			}
			if p.ShowStartDots != tt.wantStartDots {
				t.Errorf("ShowStartDots = %v, want %v", p.ShowStartDots, tt.wantStartDots)
			}
			if p.ShowEndDots != tt.wantEndDots {
				t.Errorf("ShowEndDots = %v, want %v", p.ShowEndDots, tt.wantEndDots)
			}
		})
	}
}

func intSliceEqual(a, b []int) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
