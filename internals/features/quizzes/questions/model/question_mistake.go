package model

import (
	"time"

	"github.com/google/uuid"
)

type QuestionMistakeModel struct {
	QuestionMistakeID           uint      `gorm:"column:question_mistake_id;primaryKey" json:"question_mistake_id"`
	QuestionMistakeUserID       uuid.UUID `gorm:"column:question_mistake_user_id;type:uuid;not null" json:"question_mistake_user_id"`
	QuestionMistakeSourceTypeID int       `gorm:"column:question_mistake_source_type_id;not null" json:"question_mistake_source_type_id"`
	QuestionMistakeQuestionID   uint      `gorm:"column:question_mistake_question_id;not null" json:"question_mistake_question_id"`
	CreatedAt                   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (QuestionMistakeModel) TableName() string {
	return "question_mistakes"
}
