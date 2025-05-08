package model

import (
	"time"

	"github.com/google/uuid"
)

type UserTestExam struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TestExamID      int       `gorm:"not null" json:"test_exam_id"`
	PercentageGrade int       `gorm:"not null;default:0" json:"percentage_grade"`
	TimeDuration    int       `gorm:"not null;default:0" json:"time_duration"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Opsional: Custom nama tabel jika kamu mau tetap pakai snake_case
func (UserTestExam) TableName() string {
	return "user_test_exam"
}
