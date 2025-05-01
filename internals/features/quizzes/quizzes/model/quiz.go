package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuizModel struct {
	ID            int           `json:"id" gorm:"primaryKey"`
	Name          string        `json:"name_quizzes" gorm:"type:varchar(50);unique;not null;column:name_quizzes"`
	Status        string        `json:"status" gorm:"type:varchar(10);default:pending;check:status IN ('active', 'pending', 'archived')"`
	TotalQuestion pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"total_question"`
	IconURL       string        `json:"icon_url" gorm:"type:varchar(100)"`
	CreatedAt     time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time    `json:"deleted_at" gorm:"index"`
	SectionQuizID int           `json:"section_quizzes_id" gorm:"column:section_quizzes_id"`
	CreatedBy     uuid.UUID     `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

func (QuizModel) TableName() string {
	return "quizzes"
}

// ✅ AfterSave: sinkronkan daftar quiz ke section_quizzes
func (q *QuizModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalQuizzes(tx, q.SectionQuizID)
}

// ✅ AfterDelete: sinkronkan daftar quiz ke section_quizzes
func (q *QuizModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalQuizzes(tx, q.SectionQuizID)
}

// ✅ Fungsi sinkronisasi array ID quiz ke section_quizzes.total_quizzes
func SyncTotalQuizzes(db *gorm.DB, sectionQuizID int) error {
	log.Println("[SERVICE] SyncTotalQuizzes - section_quizzes_id:", sectionQuizID)

	err := db.Exec(`
		UPDATE section_quizzes
		SET total_quizzes = (
			SELECT ARRAY_AGG(id ORDER BY id)
			FROM quizzes
			WHERE section_quizzes_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, sectionQuizID, sectionQuizID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_quizzes:", err)
	}
	return err
}
