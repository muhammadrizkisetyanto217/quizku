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

// UpdateUserUnitFromReading digunakan untuk menambahkan nilai attempt_reading pada user_unit
// ketika user menyelesaikan satu bacaan (reading) dalam unit tertentu.
//
// - Jika entry user_unit ditemukan, maka field attempt_reading akan ditambah 1.
// - Jika tidak ditemukan, tidak dilakukan create. Hanya log warning sebagai informasi.
//
// Params:
// - db: koneksi database GORM
// - userID: UUID dari user
// - unitID: ID unit tempat reading berada
//
// Return:
// - error jika terjadi kesalahan saat update database
func UpdateUserUnitFromReading(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	result := db.Model(&userUnitModel.UserUnitModel{}).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		UpdateColumn("attempt_reading", gorm.Expr("attempt_reading + 1"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Log peringatan jika user_unit tidak ditemukan (tidak auto-create)
		log.Printf("[WARNING] Tidak ditemukan user_unit untuk user_id: %s, unit_id: %d", userID, unitID)
	}
	return nil
}

// CheckAndUnsetUserUnitReadingStatus berfungsi untuk memeriksa apakah
// user masih memiliki data reading aktif pada unit tertentu.
// Jika tidak ada reading yang tercatat, maka field attempt_reading akan di-reset ke 0.
//
// Fitur ini berguna saat user menghapus semua reading, maka status attempt_reading
// juga harus dikosongkan agar progres akurat.
//
// Params:
// - db: koneksi database GORM
// - userID: UUID dari user
// - unitID: ID unit yang dicek
//
// Return:
// - error jika terjadi kegagalan query
// - nil jika berhasil, meskipun tidak ada data yang ditemukan
func CheckAndUnsetUserUnitReadingStatus(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var count int64
	err := db.Table("user_readings").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		// Reset attempt_reading jika semua reading user di unit ini telah dihapus
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("attempt_reading", 0).Error
	}

	return nil
}
