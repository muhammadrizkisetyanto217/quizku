package service

import (
	"log"

	"quizku/internals/features/users/user/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateInitialUserProfile(db *gorm.DB, userID uuid.UUID) {
	profile := models.UsersProfileModel{
		UserID: userID,
		Gender: nil, // atau models.Male jika mau default
	}
	if err := db.Create(&profile).Error; err != nil {
		log.Printf("[ERROR] Failed to create user profile: %v", err)
	}
}
