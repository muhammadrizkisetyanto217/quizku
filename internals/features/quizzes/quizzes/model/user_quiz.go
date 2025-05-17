package model

import (
	"time"

	"github.com/google/uuid"
)

type UserQuizzesModel struct {
	UserQuizID              uint      `gorm:"column:user_quiz_id;primaryKey;autoIncrement" json:"user_quiz_id"`
	UserQuizUserID          uuid.UUID `gorm:"column:user_quiz_user_id;type:uuid;not null" json:"user_quiz_user_id"`
	UserQuizQuizID          uint      `gorm:"column:user_quiz_quiz_id;not null" json:"user_quiz_quiz_id"`
	UserQuizAttempt         int       `gorm:"column:user_quiz_attempt;default:1;not null" json:"user_quiz_attempt"`
	UserQuizPercentageGrade int       `gorm:"column:user_quiz_percentage_grade;default:0;not null" json:"user_quiz_percentage_grade"`
	UserQuizTimeDuration    int       `gorm:"column:user_quiz_time_duration;default:0;not null" json:"user_quiz_time_duration"`
	UserQuizPoint           int       `gorm:"column:user_quiz_point;default:0;not null" json:"user_quiz_point"`
	CreatedAt               time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (UserQuizzesModel) TableName() string {
	return "user_quizzes"
}
