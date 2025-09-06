package repositories

import (
	"time"

	"github.com/f-alotaibi/go-starter/models"
	"gorm.io/gorm"
)

func CreatePasswordResetToken(db *gorm.DB, resetToken *models.PasswordResetToken) error {
	result := db.Create(resetToken)
	return result.Error
}

func FindPasswordResetToken(db *gorm.DB, token string) (models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	result := db.First(&resetToken, "token = ?", token)
	return resetToken, result.Error
}

func SetPasswordResetTokenAsUsed(db *gorm.DB, token string) error {
	result := db.Model(&models.PasswordResetToken{}).Where("token = ?", token).Update("used", true)
	return result.Error
}

func CleanupExpiredPasswordResetTokens(db *gorm.DB) error {
	result := db.Where("expiration < ? OR used = ?", time.Now(), true).Delete(&models.PasswordResetToken{})
	return result.Error
}
