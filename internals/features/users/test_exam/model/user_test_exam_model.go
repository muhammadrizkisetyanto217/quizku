package model

import (
	"time"

	"github.com/google/uuid"
)

type UserTestExam struct {
	UserTestExamID            uint      `gorm:"column:user_test_exam_id;primaryKey" json:"user_test_exam_id"`
	UserTestExamUserID        uuid.UUID `gorm:"column:user_test_exam_user_id;type:uuid;not null" json:"user_test_exam_user_id"`
	UserTestExamTestExamID    uint      `gorm:"column:user_test_exam_test_exam_id;not null" json:"user_test_exam_test_exam_id"`
	UserTestExamPercentageGrade int     `gorm:"column:user_test_exam_percentage_grade;default:0;not null" json:"user_test_exam_percentage_grade"`
	UserTestExamTimeDuration  int       `gorm:"column:user_test_exam_time_duration;default:0;not null" json:"user_test_exam_time_duration"`
	CreatedAt                 time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (UserTestExam) TableName() string {
	return "user_test_exams"
}