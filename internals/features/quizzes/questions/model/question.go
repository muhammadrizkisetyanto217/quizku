package model

import (
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionModel struct {
	QuestionID            uint           `gorm:"column:question_id;primaryKey" json:"question_id"`
	QuestionText          string         `gorm:"column:question_text;type:text;not null" json:"question_text"`                              // Isi pertanyaan
	QuestionAnswerChoices pq.StringArray `gorm:"column:question_answer_choices;type:text[];not null" json:"question_answer_choices"`        // Pilihan jawaban
	QuestionCorrectAnswer string         `gorm:"column:question_correct_answer;type:varchar(50);not null" json:"question_correct_answer"`   // Jawaban benar
	QuestionHelpParagraph string         `gorm:"column:question_paragraph_help;type:text;not null" json:"question_paragraph_help"`          // Paragraf bantuan jika ada
	QuestionExplanation   string         `gorm:"column:question_explanation;type:text;not null" json:"question_explanation"`                // Penjelasan mengapa jawaban benar
	QuestionAnswerText    string         `gorm:"column:question_answer_text;type:text;not null" json:"question_answer_text"`                // Ringkasan teks jawaban
	QuestionStatus        string         `gorm:"column:question_status;type:varchar(10);not null;default:'pending'" json:"question_status"` // pending, active, archived
	CreatedAt             time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (QuestionModel) TableName() string {
	return "questions"
}

// ✅ Fungsi dinamis untuk update total_question ke quizzes/evaluations/exams
func SyncTotalQuestions(db *gorm.DB, targetType int, targetID int) error {
	log.Printf("[SERVICE] SyncTotalQuestions - target_type: %d, target_id: %d\n", targetType, targetID)

	var tableName string
	switch targetType {
	case TargetTypeQuiz:
		tableName = "quizzes"
	case TargetTypeEvaluation:
		tableName = "evaluations"
	case TargetTypeExam:
		tableName = "exams"
	default:
		log.Println("[WARNING] Unknown target_type:", targetType)
		return nil
	}

	return db.Exec(`
		UPDATE `+tableName+`
		SET total_question = (
			SELECT COUNT(*)
			FROM question_links
			WHERE question_link_target_type = ? AND question_link_target_id = ?
		)
		WHERE id = ?
	`, targetType, targetID, targetID).Error
}
