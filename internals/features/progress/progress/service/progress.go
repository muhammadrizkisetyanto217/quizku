package service

import (
	"log"
	"quizku/internals/features/progress/progress/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateInitialUserProgress(db *gorm.DB, userID uuid.UUID) error {
	progress := model.UserProgress{
		UserID:      userID,
		TotalPoints: 0,
		LastUpdated: time.Now(),
	}
	if err := db.Create(&progress).Error; err != nil {
		log.Println("[ERROR] Gagal inisialisasi user_progress:", err)
		return err
	}
	log.Println("[SUCCESS] User progress berhasil diinisialisasi:", userID)
	return nil
}

func UpdateUserProgressTotal(db *gorm.DB, userID uuid.UUID) error {
	var total int64
	err := db.Table("user_point_logs").
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points), 0)").
		Scan(&total).Error
	if err != nil {
		log.Println("[ERROR] Gagal hitung total poin:", err)
		return err
	}

	// ✅ Langsung update user_progress karena pasti sudah ada
	err = db.Model(&model.UserProgress{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_points": total,
			"last_updated": time.Now(),
		}).Error

	if err != nil {
		log.Println("[ERROR] Gagal update user_progress:", err)
	}
	return err
}
