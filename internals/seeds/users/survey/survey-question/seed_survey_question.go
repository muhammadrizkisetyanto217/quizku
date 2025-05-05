package survey

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/users/survey/model"

	"gorm.io/gorm"
)

type SurveyQuestionSeed struct {
	QuestionText   string   `json:"question_text"`
	QuestionAnswer []string `json:"question_answer"`
	OrderIndex     int      `json:"order_index"`
}

func SeedSurveyQuestionsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file survey questions:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []SurveyQuestionSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	// Ambil semua question_text yang sudah ada
	var existingQuestions []string
	if err := db.Model(&model.SurveyQuestion{}).
		Select("question_text").
		Find(&existingQuestions).Error; err != nil {
		log.Fatalf("❌ Gagal ambil data existing: %v", err)
	}

	existingMap := make(map[string]bool)
	for _, q := range existingQuestions {
		existingMap[q] = true
	}

	// Filter data baru
	var newQuestions []model.SurveyQuestion
	for _, s := range seeds {
		if existingMap[s.QuestionText] {
			log.Printf("ℹ️ Pertanyaan '%s' sudah ada, dilewati.", s.QuestionText)
			continue
		}

		newQuestions = append(newQuestions, model.SurveyQuestion{
			QuestionText:   s.QuestionText,
			QuestionAnswer: s.QuestionAnswer,
			OrderIndex:     s.OrderIndex,
		})
	}

	if len(newQuestions) > 0 {
		if err := db.Create(&newQuestions).Error; err != nil {
			log.Fatalf("❌ Gagal insert survey_questions: %v", err)
		}
		log.Printf("✅ Berhasil insert %d survey questions", len(newQuestions))
	} else {
		log.Println("ℹ️ Tidak ada pertanyaan baru untuk diinsert.")
	}
}
