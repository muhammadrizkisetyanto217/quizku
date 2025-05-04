package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDailyActivity struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Date       time.Time      `gorm:"type:date;not null;index" json:"date"` // Tanggal tanpa waktu
	AmountDay  int            `gorm:"not null" json:"amount_day"`           // Rentetan hari aktif
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserDailyActivity) TableName() string {
	return "user_daily_activity"
}
