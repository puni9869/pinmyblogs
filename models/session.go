package models

import "gorm.io/gorm"

// Session is gorm sessionstore shadow  to properly delete the session
type Session struct {
	gorm.Model
	data string //lint:ignore U1000 Ignore unused function temporarily for debugging
}
