package template_functions

import (
	"strings"
	"testing"
	"time"
)

func TestDomainName(t *testing.T) {
	tests := []struct {
		link string
		want string
	}{
		{"https://www.example.com/path", "www.example.com"},
		{"http://blog.test.io", "blog.test.io"},
		{"https://example.com:8080/page", "example.com"},
		{"not-a-url", ""},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.link, func(t *testing.T) {
			got := DomainName(tt.link)
			if got != tt.want {
				t.Errorf("DomainName(%q) = %q, want %q", tt.link, got, tt.want)
			}
		})
	}
}

func TestAvatarInitials(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Alice", "A"},
		{"bob", "b"},
		{"123", "1"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := AvatarInitials(tt.input)
			if got != tt.want {
				t.Errorf("AvatarInitials(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2,3) = %d, want 5", got)
	}
	if got := Add(-1, 1); got != 0 {
		t.Errorf("Add(-1,1) = %d, want 0", got)
	}
}

func TestSub(t *testing.T) {
	if got := Sub(5, 3); got != 2 {
		t.Errorf("Sub(5,3) = %d, want 2", got)
	}
	if got := Sub(1, 5); got != -4 {
		t.Errorf("Sub(1,5) = %d, want -4", got)
	}
}

func TestAsset(t *testing.T) {
	fn := Asset("abc123")
	got := fn("app.js")
	want := "/statics/app.js?v=abc123"
	if got != want {
		t.Errorf("Asset(abc123)(app.js) = %q, want %q", got, want)
	}
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		contains string
	}{
		{"zero time", time.Time{}, "unknown time"},
		{"just now", now.Add(-10 * time.Second), "just now"},
		{"minutes ago", now.Add(-5 * time.Minute), "minutes ago"},
		{"hours ago", now.Add(-3 * time.Hour), "hours ago"},
		{"days ago", now.Add(-5 * 24 * time.Hour), "days ago"},
		{"months ago", now.Add(-60 * 24 * time.Hour), "months ago"},
		{"years ago", now.Add(-400 * 24 * time.Hour), "years ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRelativeTime(tt.input)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("FormatRelativeTime() = %q, want it to contain %q", got, tt.contains)
			}
		})
	}
}
