package models

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"
	"strings"
	"testing"
	"time"
)

func TestUserModelHidePasswordWhenLogged(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	uuid.SetRand(rnd)
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
	want := "Email:  FirstName:  Id: 67d10549-3936-4934-b1c7-e42ff60e7c81"
	got := fmt.Sprintf("%v", u)
	if strings.Compare(got, want) != 0 {
		t.Errorf("Get() = %s, want %s", got, want)
	}
}
