package models

import (
	"time"
)

type PasswordResetToken struct {
	UserID     uint      `gorm:"not null,primarykey"`
	Token      string    `gorm:"not null,unique,primarykey"`
	Expiration time.Time `gorm:"not null"`
	Used       bool      `gorm:"default:false"`
	User       User      `gorm:"foreignKey:UserID"`
}
