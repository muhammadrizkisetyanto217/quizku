package model

import (
	"time"

	"quizku/internals/features/quizzes/exams/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserExamModel struct {
	UserExamID           uint           `gorm:"column:user_exam_id;primaryKey;autoIncrement" json:"user_exam_id"`
	UserExamUserID       uuid.UUID      `gorm:"column:user_exam_user_id;not null;index:idx_user_exam_user_exam_id,priority:1;index:idx_user_exam_user_unit_id,priority:1" json:"user_exam_user_id"`
	UserExamExamID       uint           `gorm:"column:user_exam_exam_id;not null;index:idx_user_exam_user_exam_id,priority:2" json:"user_exam_exam_id"`
	UserExamUnitID       uint           `gorm:"column:user_exam_unit_id;not null;index:idx_user_exam_user_unit_id,priority:2" json:"user_exam_unit_id"`
	UserExamAttempt      int            `gorm:"column:user_exam_attempt;not null;default:1" json:"user_exam_attempt"`
	UserExamPercentageGrade int         `gorm:"column:user_exam_percentage_grade;not null;default:0" json:"user_exam_percentage_grade"`
	UserExamTimeDuration int            `gorm:"column:user_exam_time_duration;not null;default:0" json:"user_exam_time_duration"`
	UserExamPoint        int            `gorm:"column:user_exam_point;not null;default:0" json:"user_exam_point"`

	CreatedAt            time.Time      `gorm:"column:created_at;default:current_timestamp" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;default:current_timestamp" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// ✅ Table name override
func (UserExamModel) TableName() string {
	return "user_exams"
}

// ✅ Callback: Update progress after create/update/delete
func (u *UserExamModel) AfterCreate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromExam(tx, u.UserExamUserID, u.UserExamExamID, u.UserExamPercentageGrade)
}

func (u *UserExamModel) AfterUpdate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromExam(tx, u.UserExamUserID, u.UserExamExamID, u.UserExamPercentageGrade)
}

func (u *UserExamModel) AfterDelete(tx *gorm.DB) error {
	return service.CheckAndUnsetExamStatus(tx, u.UserExamUserID, u.UserExamExamID)
}
