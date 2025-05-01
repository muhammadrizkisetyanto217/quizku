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
	NameQuizzes      string `json:"name_quizzes"`
	Status           string `json:"status"`
	MaterialsQuizzes string `json:"materials_quizzes"` // optional field kalau kamu tambahkan
	IconURL          string `json:"icon_url"`
	SectionQuizID    int    `json:"section_quizzes_id"`
	CreatedBy        string `json:"created_by"`
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
		if err := db.Where("name_quizzes = ?", seed.NameQuizzes).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Quiz '%s' sudah ada, lewati...", seed.NameQuizzes)
			continue
		}

		createdByUUID := parseUUID(seed.CreatedBy)

		newQuiz := model.QuizModel{
			Name:          seed.NameQuizzes,
			Status:        seed.Status,
			IconURL:       seed.IconURL,
			SectionQuizID: seed.SectionQuizID,
			CreatedBy:     createdByUUID,
		}

		if err := db.Create(&newQuiz).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.NameQuizzes, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.NameQuizzes)
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
