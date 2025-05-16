package model

import (
	"time"

	"github.com/google/uuid"
)

type QuestionSavedModel struct {
	QuestionSavedID           uint      `gorm:"column:question_saved_id;primaryKey" json:"question_saved_id"`
	QuestionSavedUserID       uuid.UUID `gorm:"column:question_saved_user_id;type:uuid;not null" json:"question_saved_user_id"`     // Relasi ke tabel users
	QuestionSavedSourceTypeID int       `gorm:"column:question_saved_source_type_id;not null" json:"question_saved_source_type_id"` // 1 = Quiz, 2 = Evaluation, 3 = Exam
	QuestionSavedQuestionID   uint      `gorm:"column:question_saved_question_id;not null" json:"question_saved_question_id"`       // ID dari question
	CreatedAt                 time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                 // Timestamp otomatis
}

func (QuestionSavedModel) TableName() string {
	return "question_saved"
}
