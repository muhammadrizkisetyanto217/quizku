package quizzes

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizSeed struct {
	QuizName             string `json:"quiz_name"`
	QuizStatus           string `json:"quiz_status"`
	QuizIconURL          string `json:"quiz_icon_url"`
	QuizSectionQuizzesID int    `json:"quiz_section_quizzes_id"`
	QuizCreatedBy        string `json:"quiz_created_by"`
}

func SeedQuizzesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []QuizSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.QuizModel
		if err := db.Where("quiz_name = ?", seed.QuizName).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Quiz '%s' sudah ada, lewati...", seed.QuizName)
			continue
		}

		createdByUUID := parseUUID(seed.QuizCreatedBy)

		newQuiz := model.QuizModel{
			QuizName:             seed.QuizName,
			QuizStatus:           seed.QuizStatus,
			QuizIconURL:          seed.QuizIconURL,
			QuizSectionQuizzesID: seed.QuizSectionQuizzesID,
			QuizCreatedBy:        createdByUUID,
		}

		if err := db.Create(&newQuiz).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.QuizName, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.QuizName)
		}
	}
}

// helper mandiri
func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("❌ UUID tidak valid: %v", err)
	}
	return id
}
