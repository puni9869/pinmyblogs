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
