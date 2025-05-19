package model

import (
	"time"

	useSectionQuizzes "quizku/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserUnitModel struct {
	UserUnitID     uint      `gorm:"column:user_unit_id;primaryKey;autoIncrement" json:"user_unit_id"`
	UserUnitUserID uuid.UUID `gorm:"column:user_unit_user_id;type:uuid;not null;index:idx_user_unit_user_unit_unique,unique" json:"user_unit_user_id"`
	UserUnitUnitID uint      `gorm:"column:user_unit_unit_id;not null;index:idx_user_unit_user_unit_unique,unique" json:"user_unit_unit_id"`

	UserUnitAttemptReading         int            `gorm:"column:user_unit_attempt_reading;default:0;not null" json:"user_unit_attempt_reading"`
	UserUnitAttemptEvaluation      datatypes.JSON `gorm:"column:user_unit_attempt_evaluation;type:jsonb;not null;default:'{}'" json:"user_unit_attempt_evaluation"`
	UserUnitCompleteSectionQuizzes datatypes.JSON `gorm:"column:user_unit_complete_section_quizzes;type:jsonb;not null;default:'{}'" json:"user_unit_complete_section_quizzes"`

	UserUnitGradeQuiz   int `gorm:"column:user_unit_grade_quiz;default:0;check:user_unit_grade_quiz BETWEEN 0 AND 100" json:"user_unit_grade_quiz"`
	UserUnitGradeExam   int `gorm:"column:user_unit_grade_exam;default:0;check:user_unit_grade_exam BETWEEN 0 AND 100" json:"user_unit_grade_exam"`
	UserUnitGradeResult int `gorm:"column:user_unit_grade_result;default:0;check:user_unit_grade_result BETWEEN 0 AND 100" json:"user_unit_grade_result"`

	UserUnitIsPassed bool `gorm:"column:user_unit_is_passed;default:false" json:"user_unit_is_passed"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	SectionProgress []useSectionQuizzes.UserSectionQuizzesModel `gorm:"foreignKey:UserID;references:UserUnitUserID" json:"section_progress"`
}

// TableName untuk override nama tabel default
func (UserUnitModel) TableName() string {
	return "user_unit"
}
