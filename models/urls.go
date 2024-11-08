package models

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	WebLink      string
	IsActive     bool
	IsDeleted    bool
	CreatedBy    string `gorm:"index;size:255;not null"`
	Comment      string
	Summary      string
	Title        string
	Tag          string
	IsFav        bool
	IsArchived   bool
	Category     string
	IsBookMarked bool
}
