package model

import (
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	QuestionText    string         `gorm:"type:text;not null" json:"question_text"`
	QuestionAnswer  pq.StringArray `gorm:"type:text[];not null" json:"question_answer"`
	QuestionCorrect string         `gorm:"type:varchar(50);not null" json:"question_correct"`
	ParagraphHelp   string         `gorm:"type:text;not null" json:"paragraph_help"`
	ExplainQuestion string         `gorm:"type:text;not null" json:"explain_question"`
	AnswerText      string         `gorm:"type:text;not null" json:"answer_text"`
	SourceTypeID    int            `gorm:"not null" json:"source_type_id"` // 🔥 Tambahkan untuk relasi quizzes/evaluations/exams
	SourceID        uint           `gorm:"not null" json:"source_id"`      // 🔥 Tambahkan ID sumber
	DonationID      *int           `gorm:"type:int" json:"donation_id"`    // nullable
	Status          string         `gorm:"type:varchar(10);not null;default:'pending';check:status IN ('active','pending','archived')" json:"status"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (QuestionModel) TableName() string {
	return "questions"
}

// ✅ Fungsi dinamis untuk update total_question ke quizzes/evaluations/exams
func SyncTotalQuestions(db *gorm.DB, sourceTypeID int, sourceID int) error {
	log.Printf("[SERVICE] SyncTotalQuestions - source_type_id: %d, source_id: %d\n", sourceTypeID, sourceID)

	var tableName string
	switch sourceTypeID {
	case 1:
		tableName = "quizzes"
	case 2:
		tableName = "evaluations"
	case 3:
		tableName = "exams"
	default:
		log.Println("[WARNING] Unknown source_type_id:", sourceTypeID)
		return nil
	}

	return db.Exec(`
		UPDATE `+tableName+`
		SET total_question = (
			SELECT ARRAY_AGG(id ORDER BY id)
			FROM questions
			WHERE source_type_id = ? AND source_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, sourceTypeID, sourceID, sourceID).Error
}
