package models

import (
	"database/sql"

	"github.com/f-alotaibi/go-starter/models/types"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string     `gorm:"unique;not null"`
	Email             string     `gorm:"unique;not null"`
	Password          []byte     `gorm:"not null"`
	Role              types.Role `gorm:"not null;default:user"`
	LastPasswordReset sql.NullTime
}
