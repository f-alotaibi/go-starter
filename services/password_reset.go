package services

import (
	"log"
	"time"

	"github.com/f-alotaibi/go-starter/repositories"
	"gorm.io/gorm"
)

func StartPasswordResetTokenCleanup(db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	go func(db *gorm.DB) {
		for range ticker.C {
			if err := repositories.CleanupExpiredPasswordResetTokens(db); err != nil {
				log.Println("Failed to cleanup tokens:", err)
			} else {
				log.Println("Expired tokens cleaned up")
			}
		}
	}(db)
}
