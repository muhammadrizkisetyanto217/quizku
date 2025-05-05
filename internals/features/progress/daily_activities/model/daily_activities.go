package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDailyActivity struct {
	ID           uint           `gorm:"primaryKey" json:"id"`                                              // SERIAL
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`                           // UUID
	Date         time.Time      `gorm:"type:date;not null" json:"date"`                                    // tanggal logika (sama dgn activity_date)
	ActivityDate time.Time      `gorm:"type:date;not null;uniqueIndex:idx_user_date" json:"activity_date"` // tanggal aktivitas
	AmountDay    int            `gorm:"not null;default:1" json:"amount_day"`                              // jumlah hari aktif
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (UserDailyActivity) TableName() string {
	return "user_daily_activity"
}
