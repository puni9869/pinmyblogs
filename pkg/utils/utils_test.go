package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHashifySHA256(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := HashifySHA256(tt.input)
			if got != tt.want {
				t.Errorf("HashifySHA256(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHashifyAndCompareBCrypt(t *testing.T) {
	password := "mysecretpassword"
	hashed := HashifyBCrypt(password)

	if err := CompareBCrypt(hashed, password); err != nil {
		t.Errorf("CompareBCrypt should match: %v", err)
	}

	if err := CompareBCrypt(hashed, "wrongpassword"); err == nil {
		t.Error("CompareBCrypt should fail for wrong password")
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"test+tag@domain.co", true},
		{"invalid", false},
		{"@missing.com", false},
		{"", false},
		{"no@domain", true}, // mail.ParseAddress accepts this
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.valid {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}

func TestGetPageAndLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		query     string
		wantPage  int
		wantLimit int
	}{
		{"defaults", "", 1, 50},
		{"custom values", "page=3&limit=20", 3, 20},
		{"negative page defaults to 1", "page=-1&limit=10", 1, 10},
		{"zero page defaults to 1", "page=0", 1, 50},
		{"limit capped at 100", "page=1&limit=200", 1, 100},
		{"negative limit defaults", "limit=-5", 1, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			page, limit := GetPageAndLimit(c)
			if page != tt.wantPage {
				t.Errorf("page = %d, want %d", page, tt.wantPage)
			}
			if limit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", limit, tt.wantLimit)
			}
		})
	}
}
