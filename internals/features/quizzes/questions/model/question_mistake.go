package model

import (
	"time"

	"github.com/google/uuid"
)

type QuestionMistakeModel struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`           // Relasi ke tabel users
	SourceTypeID int       `gorm:"not null" json:"source_type_id"`              // 1 = Quiz, 2 = Evaluation, 3 = Exam
	QuestionID   uint      `gorm:"not null" json:"question_id"`                 // ID dari question
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Timestamp otomatis
}

// TableName untuk mapping ke tabel database
func (QuestionMistakeModel) TableName() string {
	return "question_mistakes"
}
