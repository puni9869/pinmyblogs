package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"primaryKey;unique;type:uuid;default:gen_random_uuid()"` // Standard field for the primary key
	FirstName       string    // A regular string field
	LastName        string    // A regular string field
	DisplayName     string    // public displayname
	Password        string    // Password is a type of hash
	EmailVerifyHash string    // EmailVerifyHash verification hash
	IsEmailVerified bool      // Verify the email for first time
	IsActive        bool      // Active or InActive means temporary disable account
	IsProfilePublic bool      // Profile is public or private
	Email           string    `gorm:"unique;not null"` // A pointer to a string, allowing for null values
	ActivatedAt     time.Time // Uses time.Time for nullable time fields
	CreatedAt       time.Time // Automatically managed by GORM for creation time
	UpdatedAt       time.Time // Automatically managed by GORM for update time
}
