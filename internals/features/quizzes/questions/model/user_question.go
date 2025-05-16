package model

import (
	"time"

	"github.com/google/uuid"
)

type UserQuestionModel struct {
	UserQuestionID             uint      `gorm:"column:user_question_id;primaryKey" json:"user_question_id"`
	UserQuestionUserID         uuid.UUID `gorm:"column:user_question_user_id;type:uuid;not null" json:"user_question_user_id"`
	UserQuestionQuestionID     int       `gorm:"column:user_question_question_id;not null" json:"user_question_question_id"`
	UserQuestionSelectedAnswer string    `gorm:"column:user_question_selected_answer;type:text;not null" json:"user_question_selected_answer"`
	UserQuestionIsCorrect      bool      `gorm:"column:user_question_is_correct;not null" json:"user_question_is_correct"`
	UserQuestionSourceTypeID   int       `gorm:"column:user_question_source_type_id;not null" json:"user_question_source_type_id"`
	UserQuestionSourceID       int       `gorm:"column:user_question_source_id;not null" json:"user_question_source_id"`
	CreatedAt                  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (UserQuestionModel) TableName() string {
	return "user_questions"
}
