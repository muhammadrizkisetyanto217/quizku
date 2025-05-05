package model

import (
	"time"

	"github.com/lib/pq"
)

type SurveyQuestion struct {
	ID             int            `gorm:"primaryKey" json:"id"`
	QuestionText   string         `gorm:"type:text;not null" json:"question_text"`
	QuestionAnswer pq.StringArray `gorm:"type:text[]" json:"question_answer,omitempty"`
	OrderIndex     int            `gorm:"not null;index" json:"order_index"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyQuestion) TableName() string {
	return "survey_questions"
}
