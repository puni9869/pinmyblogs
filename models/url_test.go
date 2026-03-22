package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUrlJSONMarshal_ExcludesHiddenFields(t *testing.T) {
	u := Url{
		ID:        1,
		WebLink:   "https://example.com",
		IsActive:  true,
		IsDeleted: false,
		CreatedBy: "user@test.com",
		Comment:   "secret comment",
		Summary:   "A summary",
		Title:     "Example",
		Tag:       "tech",
		IsFav:     true,
		Category:  "hidden-category",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// These fields have json:"-" and must NOT appear
	hiddenFields := []string{"DeletedAt", "IsActive", "Comment", "Category"}
	for _, f := range hiddenFields {
		if _, ok := m[f]; ok {
			t.Errorf("field %q should be excluded from JSON but was present", f)
		}
	}

	// These fields must appear
	visibleFields := []string{"id", "webLink", "isDeleted", "createdBy", "summary", "title", "tag", "isFav", "isArchived", "isBookMarked"}
	for _, f := range visibleFields {
		if _, ok := m[f]; !ok {
			t.Errorf("field %q should be in JSON but was missing", f)
		}
	}
}

func TestUrlJSONMarshal_Values(t *testing.T) {
	u := Url{
		ID:         42,
		WebLink:    "https://go.dev",
		Title:      "Go",
		Tag:        "lang",
		IsFav:      true,
		IsArchived: false,
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if m["id"].(float64) != 42 {
		t.Errorf("id = %v, want 42", m["id"])
	}
	if m["webLink"] != "https://go.dev" {
		t.Errorf("webLink = %v, want https://go.dev", m["webLink"])
	}
	if m["isFav"] != true {
		t.Errorf("isFav = %v, want true", m["isFav"])
	}
}

func TestUrlZeroValues(t *testing.T) {
	u := Url{}

	if u.ID != 0 {
		t.Errorf("zero Url.ID = %d, want 0", u.ID)
	}
	if u.IsActive {
		t.Error("zero Url.IsActive should be false")
	}
	if u.IsFav {
		t.Error("zero Url.IsFav should be false")
	}
	if u.IsArchived {
		t.Error("zero Url.IsArchived should be false")
	}
	if u.IsDeleted {
		t.Error("zero Url.IsDeleted should be false")
	}
	if u.IsBookMarked {
		t.Error("zero Url.IsBookMarked should be false")
	}
	if u.WebLink != "" {
		t.Errorf("zero Url.WebLink = %q, want empty", u.WebLink)
	}
}

func TestUrlJSONRoundTrip(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	original := Url{
		ID:           10,
		CreatedAt:    now,
		UpdatedAt:    now,
		WebLink:      "https://example.com/article",
		IsDeleted:    false,
		CreatedBy:    "alice@test.com",
		Summary:      "Great article",
		Title:        "Test Article",
		Tag:          "reading",
		IsFav:        true,
		IsArchived:   false,
		IsBookMarked: true,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Url
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}
	if decoded.WebLink != original.WebLink {
		t.Errorf("WebLink = %q, want %q", decoded.WebLink, original.WebLink)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if decoded.IsFav != original.IsFav {
		t.Errorf("IsFav = %v, want %v", decoded.IsFav, original.IsFav)
	}
	if decoded.IsBookMarked != original.IsBookMarked {
		t.Errorf("IsBookMarked = %v, want %v", decoded.IsBookMarked, original.IsBookMarked)
	}
	// Hidden fields should be zero after round-trip
	if decoded.IsActive {
		t.Error("IsActive should be false after JSON round-trip (json:\"-\")")
	}
	if decoded.Comment != "" {
		t.Errorf("Comment should be empty after JSON round-trip, got %q", decoded.Comment)
	}
	if decoded.Category != "" {
		t.Errorf("Category should be empty after JSON round-trip, got %q", decoded.Category)
	}
}
