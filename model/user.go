package model

type User struct {
	Default
	FullName *string `json:"full_name" gorm:"not null;"`
	UserName *string `json:"user_name" gorm:"not null;"`
	Password *string `json:"password" gorm:"not null;"`
}
