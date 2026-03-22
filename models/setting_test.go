package models

import (
	"testing"
)

func TestSettingZeroValues(t *testing.T) {
	s := Setting{}

	if s.CreatedBy != "" {
		t.Errorf("zero Setting.CreatedBy = %q, want empty", s.CreatedBy)
	}
	if s.Action != "" {
		t.Errorf("zero Setting.Action = %q, want empty", s.Action)
	}
	if s.Value != "" {
		t.Errorf("zero Setting.Value = %q, want empty", s.Value)
	}
	if s.IsShowCount {
		t.Error("zero Setting.IsShowCount should be false")
	}
	if s.Categories != nil {
		t.Errorf("zero Setting.Categories should be nil, got %v", s.Categories)
	}
}

func TestSettingWithValues(t *testing.T) {
	s := Setting{
		CreatedBy:   "user@test.com",
		Action:      "sideNav",
		Value:       "hide",
		IsShowCount: true,
		Categories:  []string{"tech", "news"},
	}

	if s.CreatedBy != "user@test.com" {
		t.Errorf("CreatedBy = %q, want %q", s.CreatedBy, "user@test.com")
	}
	if s.Action != "sideNav" {
		t.Errorf("Action = %q, want %q", s.Action, "sideNav")
	}
	if s.Value != "hide" {
		t.Errorf("Value = %q, want %q", s.Value, "hide")
	}
	if !s.IsShowCount {
		t.Error("IsShowCount should be true")
	}
	if len(s.Categories) != 2 {
		t.Errorf("Categories len = %d, want 2", len(s.Categories))
	}
}
