package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserSectionQuizzesModel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	SectionQuizzesID uint           `gorm:"not null" json:"section_quizzes_id"`
	CompleteQuiz     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"complete_quiz"`
	TotalQuiz        pq.Int64Array  `gorm:"type:integer[];not null;default:'{}'" json:"total_quiz"`
	GradeResult      int            `gorm:"default:0" json:"grade_result"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (UserSectionQuizzesModel) TableName() string {
	return "user_section_quizzes"
}
