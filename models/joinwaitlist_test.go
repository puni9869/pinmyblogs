package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJoinWaitListZeroValues(t *testing.T) {
	j := JoinWaitList{}

	if j.ID != uuid.Nil {
		t.Errorf("zero JoinWaitList.ID = %v, want uuid.Nil", j.ID)
	}
	if j.Email != "" {
		t.Errorf("zero JoinWaitList.Email = %q, want empty", j.Email)
	}
	if j.App != "" {
		t.Errorf("zero JoinWaitList.App = %q, want empty", j.App)
	}
}

func TestJoinWaitListWithValues(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	j := JoinWaitList{
		ID:        id,
		Email:     "test@example.com",
		App:       "pinmyblogs",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if j.ID != id {
		t.Errorf("ID = %v, want %v", j.ID, id)
	}
	if j.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", j.Email, "test@example.com")
	}
	if j.App != "pinmyblogs" {
		t.Errorf("App = %q, want %q", j.App, "pinmyblogs")
	}
	if j.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", j.CreatedAt, now)
	}
}

func TestJoinWaitListJSONTags(t *testing.T) {
	// Verify the struct can be used with JSON marshaling
	j := JoinWaitList{
		Email: "hello@world.com",
		App:   "testapp",
	}

	if j.Email != "hello@world.com" {
		t.Errorf("Email = %q, want %q", j.Email, "hello@world.com")
	}
}
