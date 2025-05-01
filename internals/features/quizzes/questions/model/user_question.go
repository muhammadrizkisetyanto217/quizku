package model

import (
	"time"

	"github.com/google/uuid"
)

type UserQuestionModel struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	QuestionID     uint      `gorm:"not null" json:"question_id"`
	SelectedAnswer string    `gorm:"type:text;not null" json:"selected_answer"`
	IsCorrect      bool      `gorm:"not null" json:"is_correct"`
	SourceTypeID   int       `gorm:"not null" json:"source_type_id"` // 1 = Quiz, 2 = Evaluation, 3 = Exam
	SourceID       uint      `gorm:"not null" json:"source_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (UserQuestionModel) TableName() string {
	return "user_questions"
}
