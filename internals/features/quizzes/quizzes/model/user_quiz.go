package model

import (
	"time"

	"github.com/google/uuid"
)

type UserQuizzesModel struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	QuizID          uint      `gorm:"column:quiz_id;not null" json:"quiz_id"` 
	Attempt         int       `gorm:"default:1;not null" json:"attempt"`
	PercentageGrade int       `gorm:"default:0;not null" json:"percentage_grade"`
	TimeDuration    int       `gorm:"default:0;not null" json:"time_duration"`
	Point           int       `gorm:"default:0;not null" json:"point"`
	CreatedAt       time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

// TableName memastikan GORM menggunakan tabel "user_quizzes"
func (UserQuizzesModel) TableName() string {
	return "user_quizzes"
}