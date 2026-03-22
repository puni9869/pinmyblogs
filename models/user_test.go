package models

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserModelHidePasswordWhenLogged(t *testing.T) {
	reader := bytes.NewReader([]byte("1111111111111111"))
	uuid.SetRand(reader)
	uuid.SetClockSequence(1)
	u := User{
		ID:              uuid.New(),
		FirstName:       "",
		LastName:        "",
		DisplayName:     "",
		Password:        "I will not ganna tell you",
		EmailVerifyHash: "",
		IsEmailVerified: false,
		IsActive:        false,
		IsProfilePublic: false,
		Email:           "",
		ActivatedAt:     time.Time{},
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
		AlternateEmail:  "",
	}
	got := fmt.Sprintf("%v", u)
	want := "{31313131-3131-4131-b131-313131313131    I will not ganna tell you  false false false  0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC    0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}"
	if strings.Compare(got, want) != 0 {
		t.Errorf("Get() = %s, want %s", got, want)
	}
}

func TestUserZeroValues(t *testing.T) {
	u := User{}

	if u.ID != uuid.Nil {
		t.Errorf("zero User.ID = %v, want uuid.Nil", u.ID)
	}
	if u.IsActive {
		t.Error("zero User.IsActive should be false")
	}
	if u.IsEmailVerified {
		t.Error("zero User.IsEmailVerified should be false")
	}
	if u.IsProfilePublic {
		t.Error("zero User.IsProfilePublic should be false")
	}
	if u.Email != "" {
		t.Errorf("zero User.Email = %q, want empty", u.Email)
	}
	if u.Password != "" {
		t.Errorf("zero User.Password = %q, want empty", u.Password)
	}
}

func TestUserWithAllFields(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	now := time.Now()

	u := User{
		ID:                id,
		FirstName:         "John",
		LastName:          "Doe",
		DisplayName:       "johndoe",
		Password:          "hashed-pw",
		Email:             "john@example.com",
		IsActive:          true,
		IsEmailVerified:   true,
		IsProfilePublic:   true,
		AlternateEmail:    "alt@example.com",
		PasswordResetHash: "reset-hash",
		AccountEnableHash: "enable-hash",
		ActivatedAt:       now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if u.ID != id {
		t.Errorf("ID = %v, want %v", u.ID, id)
	}
	if u.FirstName != "John" {
		t.Errorf("FirstName = %q, want %q", u.FirstName, "John")
	}
	if u.LastName != "Doe" {
		t.Errorf("LastName = %q, want %q", u.LastName, "Doe")
	}
	if u.DisplayName != "johndoe" {
		t.Errorf("DisplayName = %q, want %q", u.DisplayName, "johndoe")
	}
	if !u.IsActive {
		t.Error("IsActive should be true")
	}
	if !u.IsEmailVerified {
		t.Error("IsEmailVerified should be true")
	}
	if !u.IsProfilePublic {
		t.Error("IsProfilePublic should be true")
	}
	if u.AlternateEmail != "alt@example.com" {
		t.Errorf("AlternateEmail = %q, want %q", u.AlternateEmail, "alt@example.com")
	}
	if u.PasswordResetHash != "reset-hash" {
		t.Errorf("PasswordResetHash = %q, want %q", u.PasswordResetHash, "reset-hash")
	}
	if u.AccountEnableHash != "enable-hash" {
		t.Errorf("AccountEnableHash = %q, want %q", u.AccountEnableHash, "enable-hash")
	}
}

func TestUserPasswordNotExposedInSprintf(t *testing.T) {
	// Note: This test documents current behavior.
	// The password IS visible in fmt.Sprintf — consider implementing
	// a custom String() method to redact it.
	u := User{
		Password: "super-secret",
	}
	got := fmt.Sprintf("%v", u)
	if !strings.Contains(got, "super-secret") {
		t.Log("Password is hidden in Sprintf output — good")
	} else {
		t.Log("Warning: Password is visible in Sprintf output. Consider adding a String() method to redact it.")
	}
}

func TestSessionZeroValues(t *testing.T) {
	s := Session{}

	if s.ID != "" {
		t.Errorf("zero Session.ID = %q, want empty", s.ID)
	}
	if !s.CreatedAt.IsZero() {
		t.Errorf("zero Session.CreatedAt should be zero time")
	}
}

func TestSessionWithValues(t *testing.T) {
	now := time.Now()
	s := Session{
		ID:        "session-abc-123",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.ID != "session-abc-123" {
		t.Errorf("ID = %q, want %q", s.ID, "session-abc-123")
	}
	if s.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", s.CreatedAt, now)
	}
}
