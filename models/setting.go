package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	CreatedBy   string         `gorm:"index;size:255;not null"`
	Categories  pq.StringArray `gorm:"type:text[]"`
	IsShowCount bool
}
