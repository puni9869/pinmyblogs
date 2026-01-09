package models

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model

	CreatedBy string `gorm:"size:255;not null;uniqueIndex:idx_user_action"`
	Action    string `gorm:"size:255;not null;uniqueIndex:idx_user_action"`

	Categories  []string `gorm:"type:text[]"`
	IsShowCount bool
	Value       string `gorm:"size:255;not null"`
}
