package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password []byte `gorm:"not null"`
	Role     Role   `gorm:"not null;default:0"`
}
