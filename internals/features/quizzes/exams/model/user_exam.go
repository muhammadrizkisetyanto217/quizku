package model

import (
	"time"

	"quizku/internals/features/quizzes/exams/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserExamModel struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uuid.UUID      `gorm:"not null;index:idx_user_exams_user_id_exam_id,priority:1;index:idx_user_exams_user_id_unit_id,priority:1" json:"user_id"`
	ExamID          uint           `gorm:"not null;index:idx_user_exams_user_id_exam_id,priority:2" json:"exam_id"`
	UnitID          uint           `gorm:"not null;index:idx_user_exams_user_id_unit_id,priority:2" json:"unit_id"`
	Attempt         int            `gorm:"default:1;not null" json:"attempt"`
	PercentageGrade int            `gorm:"default:0;not null" json:"percentage_grade"`
	TimeDuration    int            `gorm:"default:0;not null" json:"time_duration"`
	Point           int            `gorm:"default:0;not null" json:"point"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}




func (UserExamModel) TableName() string {
	return "user_exams"
}
func (u *UserExamModel) AfterCreate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromExam(tx, u.UserID, u.ExamID, u.PercentageGrade)
}

func (u *UserExamModel) AfterUpdate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromExam(tx, u.UserID, u.ExamID, u.PercentageGrade)
}

func (u *UserExamModel) AfterDelete(tx *gorm.DB) error {
	return service.CheckAndUnsetExamStatus(tx, u.UserID, u.ExamID)
}