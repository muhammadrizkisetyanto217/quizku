package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SourceTypeID    int            `gorm:"not null" json:"source_type_id"` // 1=Quiz, 2=Evaluation, 3=Exam
	SourceID        uint           `gorm:"not null" json:"source_id"`      // quizzes_id / evaluation_id / exam_id
	QuestionText    string         `gorm:"type:varchar(200);not null" json:"question_text"`
	QuestionAnswer  pq.StringArray `gorm:"type:text[];not null" json:"question_answer"`
	QuestionCorrect string         `gorm:"type:varchar(50);not null" json:"question_correct"`
	TooltipsID      pq.Int64Array  `gorm:"type:int[]" json:"tooltips_id,omitempty"` // hanya digunakan jika source_type_id = 1
	Status          string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	ParagraphHelp   string         `gorm:"type:text;not null" json:"paragraph_help"`
	ExplainQuestion string         `gorm:"type:text;not null" json:"explain_question"`
	AnswerText      string         `gorm:"type:text;not null" json:"answer_text"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName mengatur nama tabel sesuai struktur di database
func (QuestionModel) TableName() string {
	return "questions"
}

// MarshalJSON untuk menyesuaikan array Tooltips agar bisa terbaca JSON
func (q QuestionModel) MarshalJSON() ([]byte, error) {
	type Alias QuestionModel
	return json.Marshal(&struct {
		TooltipsID []int64 `json:"tooltips_id"`
		*Alias
	}{
		TooltipsID: []int64(q.TooltipsID),
		Alias:      (*Alias)(&q),
	})
}

func (q *QuestionModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalQuestions(tx, int(q.SourceTypeID), int(q.SourceID))
}

func (q *QuestionModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalQuestions(tx, int(q.SourceTypeID), int(q.SourceID))
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
