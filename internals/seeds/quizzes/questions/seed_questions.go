package question

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/quizzes/questions/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionSeed struct {
	SourceTypeID    int      `json:"source_type_id"`
	SourceID        uint     `json:"source_id"`
	QuestionText    string   `json:"question_text"`
	QuestionAnswer  []string `json:"question_answer"`
	QuestionCorrect string   `json:"question_correct"`
	ParagraphHelp   string   `json:"paragraph_help"`
	ExplainQuestion string   `json:"explain_question"`
	AnswerText      string   `json:"answer_text"`
}

func SeedQuestionsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []QuestionSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.QuestionModel
		if err := db.Where("question_text = ? AND source_id = ? AND source_type_id = ?", seed.QuestionText, seed.SourceID, seed.SourceTypeID).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Soal '%s' sudah ada, lewati...", seed.QuestionText)
			continue
		}

		question := model.QuestionModel{
			QuestionText:    seed.QuestionText,
			QuestionAnswer:  pq.StringArray(seed.QuestionAnswer),
			QuestionCorrect: seed.QuestionCorrect,
			Status:          "active",
			ParagraphHelp:   seed.ParagraphHelp,
			ExplainQuestion: seed.ExplainQuestion,
			AnswerText:      seed.AnswerText,
		}

		if err := db.Create(&question).Error; err != nil {
			log.Printf("❌ Gagal insert soal '%s': %v", seed.QuestionText, err)
		} else {
			log.Printf("✅ Berhasil insert soal '%s'", seed.QuestionText)
		}
	}
}
