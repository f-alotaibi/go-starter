package repositories

import (
	"github.com/f-alotaibi/go-starter/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	result := db.Create(user)
	return result.Error
}

func FindUser(db *gorm.DB, username string, email string) (models.User, error) {
	var user models.User
	result := db.First(&user, "username = ? OR email = ?", username, email)
	return user, result.Error
}

func FindUserByUsername(db *gorm.DB, username string) (models.User, error) {
	var user models.User
	result := db.First(&user, "username = ?", username)
	return user, result.Error
}
