package repositories

import (
	"database/sql"
	"time"

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

func FindUserByEmail(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	result := db.First(&user, "email = ?", email)
	return user, result.Error
}

func UpdateUserPassword(db *gorm.DB, id uint, hashedPassword []byte) error {
	result := db.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Password: hashedPassword, LastPasswordReset: sql.NullTime{Time: time.Now(), Valid: true}})
	return result.Error
}
