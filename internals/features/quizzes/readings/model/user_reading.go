package model

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
	// "quiz-fiber/internals/features/quizzes/reading/service"
)

type UserReading struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"not null;index:idx_user_readings_user_id_reading_id" json:"user_id"`
	ReadingID uint      `gorm:"not null;index:idx_user_readings_user_id_reading_id" json:"reading_id"`
	UnitID    uint      `gorm:"not null;index:idx_user_readings_user_id_unit_id" json:"unit_id"`
	Attempt   int       `gorm:"default:1;not null" json:"attempt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserReading) TableName() string {
	return "user_readings"
}
