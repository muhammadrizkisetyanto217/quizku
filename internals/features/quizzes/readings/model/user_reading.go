package model

import (
	"time"

	"github.com/google/uuid"
)

type UserReading struct {
	UserReadingID        uint      `gorm:"column:user_reading_id;primaryKey;autoIncrement" json:"user_reading_id"`
	UserReadingUserID    uuid.UUID `gorm:"column:user_reading_user_id;type:uuid;not null;index:idx_user_readings_user_id_reading_id,priority:1;index:idx_user_readings_user_id_unit_id,priority:1" json:"user_reading_user_id"`
	UserReadingReadingID uint      `gorm:"column:user_reading_reading_id;not null;index:idx_user_readings_user_id_reading_id,priority:2" json:"user_reading_reading_id"`
	UserReadingUnitID    uint      `gorm:"column:user_reading_unit_id;not null;index:idx_user_readings_user_id_unit_id,priority:2" json:"user_reading_unit_id"`
	UserReadingAttempt   int       `gorm:"column:user_reading_attempt;default:1;not null" json:"user_reading_attempt"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ✅ Mapping nama tabel
func (UserReading) TableName() string {
	return "user_readings"
}
