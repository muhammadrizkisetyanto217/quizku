package model

import (
	"time"

	"gorm.io/gorm"
)

type DonationQuestionModel struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	DonationID     uint           `gorm:"not null" json:"donation_id"`
	QuestionID     uint           `gorm:"not null" json:"question_id"`
	UserProgressID *uint          `gorm:"type:int" json:"user_progress_id"` // nullable
	UserMessage    string         `gorm:"type:text" json:"user_message"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (DonationQuestionModel) TableName() string {
	return "donation_questions"
}
