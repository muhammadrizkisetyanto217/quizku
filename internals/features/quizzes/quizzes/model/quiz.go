package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuizModel struct {
	QuizID               int           `gorm:"column:quiz_id;primaryKey;autoIncrement" json:"quiz_id"`
	QuizName             string        `gorm:"column:quiz_name;type:varchar(50);unique;not null" json:"quiz_name"`
	QuizStatus           string        `gorm:"column:quiz_status;type:varchar(10);default:'pending';check:quiz_status IN ('active', 'pending', 'archived')" json:"quiz_status"`
	QuizTotalQuestion    pq.Int64Array `gorm:"column:quiz_total_question;type:integer[];default:'{}'" json:"quiz_total_question"`
	QuizIconURL          string        `gorm:"column:quiz_icon_url;type:varchar(100)" json:"quiz_icon_url"`
	CreatedAt            time.Time     `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt            *time.Time    `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	QuizSectionQuizzesID int           `gorm:"column:quiz_section_quizzes_id;not null" json:"quiz_section_quizzes_id"`
	QuizCreatedBy        uuid.UUID     `gorm:"column:quiz_created_by;type:uuid;not null;constraint:OnDelete:CASCADE" json:"quiz_created_by"`
}

func (QuizModel) TableName() string {
	return "quizzes"
}
func (q *QuizModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalQuizzes(tx, q.QuizSectionQuizzesID)
}




func (q *QuizModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalQuizzes(tx, q.QuizSectionQuizzesID)
}

func SyncTotalQuizzes(db *gorm.DB, sectionQuizID int) error {
	log.Println("[SERVICE] SyncTotalQuizzes - quiz_section_quizzes_id:", sectionQuizID)

	err := db.Exec(`
		UPDATE section_quizzes
		SET section_quizzes_total_quizzes = (
			SELECT ARRAY_AGG(quiz_id ORDER BY quiz_id)
			FROM quizzes
			WHERE quiz_section_quizzes_id = ? AND deleted_at IS NULL
		)
		WHERE section_quizzes_id = ?
	`, sectionQuizID, sectionQuizID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_quizzes:", err)
	}
	return err
}
