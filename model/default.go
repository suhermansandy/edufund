package model

import "time"

// Default model database
type Default struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}
