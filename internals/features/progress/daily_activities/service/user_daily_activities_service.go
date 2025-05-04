package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"quizku/internals/features/progress/daily_activities/model"
)

func UpdateOrInsertDailyActivity(db *gorm.DB, userID uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)
	var existing model.UserDailyActivity

	// Cek apakah sudah ada aktivitas hari ini
	err := db.Where("user_id = ? AND date = ?", userID, today).First(&existing).Error
	if err == nil {
		// Sudah ada: update timestamp
		return db.Model(&existing).Update("updated_at", time.Now()).Error
	}

	// Ambil aktivitas terakhir user (jika ada)
	var lastActivity model.UserDailyActivity
	err = db.Where("user_id = ?", userID).Order("date DESC").First(&lastActivity).Error

	var newAmountDay int
	if err == nil && lastActivity.Date.Add(24*time.Hour).Equal(today) {
		// Lanjutan dari kemarin
		newAmountDay = lastActivity.AmountDay + 1
	} else {
		// Tidak aktif atau hari pertama
		newAmountDay = 1
	}

	newActivity := model.UserDailyActivity{
		UserID:    userID,
		Date:      today,
		AmountDay: newAmountDay,
	}

	return db.Create(&newActivity).Error
}
