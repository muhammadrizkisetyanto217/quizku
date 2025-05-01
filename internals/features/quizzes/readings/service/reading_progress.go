package service

import (
	"log"
	userUnitModel "quizku/internals/features/lessons/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//////////////////////////////////////////////////////////
// === BAGIAN UNTUK USER READING ===
//////////////////////////////////////////////////////////

func UpdateUserUnitFromReading(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	// Langsung update, tanpa create jika tidak ditemukan
	result := db.Model(&userUnitModel.UserUnitModel{}).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		UpdateColumn("attempt_reading", gorm.Expr("attempt_reading + 1"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Safety: log jika tidak ditemukan (tidak membuat)
		log.Printf("[WARNING] Tidak ditemukan user_unit untuk user_id: %s, unit_id: %d", userID, unitID)
	}
	return nil
}

func CheckAndUnsetUserUnitReadingStatus(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var count int64
	err := db.Table("user_readings").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("attempt_reading", 0).Error
	}

	return nil
}
