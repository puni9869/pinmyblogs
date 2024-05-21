package models

import (
	"time"

	"gorm.io/gorm"
)

// Session is gorm sessionstore shadow  to properly delete the session
type Session struct {
	ID        string `gorm:"primarykey;size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	data      string         //lint:ignore U1000 Ignore unused function temporarily for debugging
}
