package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type DeletedAt sql.NullTime
type Url struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	WebLink      string         `json:"webLink"`
	IsActive     bool           `json:"-"`
	IsDeleted    bool           `json:"isDeleted"`
	CreatedBy    string         `gorm:"index;size:255;not null" json:"createdBy"`
	Comment      string         `json:"-"`
	Summary      string         `json:"summary"`
	Title        string         `json:"title"`
	Tag          string         `json:"tag"`
	IsFav        bool           `json:"isFav"`
	IsArchived   bool           `json:"isArchived"`
	Category     string         `json:"-"`
	IsBookMarked bool           `json:"isBookMarked"`
}
