package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSectionQuizzesModel struct {
	UserSectionQuizzesID               uint           `gorm:"column:user_section_quizzes_id;primaryKey;autoIncrement" json:"user_section_quizzes_id"`
	UserSectionQuizzesUserID           uuid.UUID      `gorm:"column:user_section_quizzes_user_id;type:uuid;not null" json:"user_section_quizzes_user_id"`
	UserSectionQuizzesSectionQuizzesID uint           `gorm:"column:user_section_quizzes_section_quizzes_id;not null" json:"user_section_quizzes_section_quizzes_id"`
	UserSectionQuizzesCompleteQuiz     datatypes.JSON `gorm:"column:user_section_quizzes_complete_quiz;type:jsonb;not null;default:'{}'" json:"user_section_quizzes_complete_quiz"`
	UserSectionQuizzesGradeResult      int            `gorm:"column:user_section_quizzes_grade_result;not null;default:0" json:"user_section_quizzes_grade_result"`
	CreatedAt                          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt                          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ✅ Mapping nama tabel
func (UserSectionQuizzesModel) TableName() string {
	return "user_section_quizzes"
}
