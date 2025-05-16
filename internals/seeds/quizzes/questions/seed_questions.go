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
	QuestionText          string   `json:"question_text"`
	QuestionAnswerChoices []string `json:"question_answer_choices"`
	QuestionCorrectAnswer string   `json:"question_correct_answer"`
	QuestionHelpParagraph string   `json:"question_help_paragraph"`
	QuestionExplanation   string   `json:"question_explanation"`
	QuestionAnswerText    string   `json:"question_answer_text"`
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
		// Cek duplikat berdasarkan QuestionText
		var existing model.QuestionModel
		if err := db.Where("question_text = ?", seed.QuestionText).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Soal '%s' sudah ada, lewati...", seed.QuestionText)
			continue
		}

		question := model.QuestionModel{
			QuestionText:          seed.QuestionText,
			QuestionAnswerChoices: pq.StringArray(seed.QuestionAnswerChoices),
			QuestionCorrectAnswer: seed.QuestionCorrectAnswer,
			QuestionHelpParagraph: seed.QuestionHelpParagraph,
			QuestionExplanation:   seed.QuestionExplanation,
			QuestionAnswerText:    seed.QuestionAnswerText,
			QuestionStatus:        "active",
		}

		if err := db.Create(&question).Error; err != nil {
			log.Printf("❌ Gagal insert soal '%s': %v", seed.QuestionText, err)
		} else {
			log.Printf("✅ Berhasil insert soal '%s'", seed.QuestionText)
		}
	}
}
