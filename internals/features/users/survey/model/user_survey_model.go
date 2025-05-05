package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSurvey struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SurveyQuestionID int       `gorm:"not null;index" json:"survey_question_id"`
	UserAnswer       string    `gorm:"type:text;not null" json:"user_answer"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (UserSurvey) TableName() string {
	return "user_surveys"
}
