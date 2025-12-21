package models

import (
	"time"

	"github.com/google/uuid"
)

type JoinWaitList struct {
	ID        uuid.UUID `gorm:"primaryKey;unique;type:uuid;default:(gen_random_uuid())"` // Standard field for the primary key
	Email     string    `gorm:"index;size:255" json:"email"`                             // A pointer to a string, allowing for null values
	App       string    `gorm:"index;size:255" json:"app"`
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}
