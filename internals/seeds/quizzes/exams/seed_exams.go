package exam

import (
	"encoding/json"
	"log"
	"os"

	"quizku/internals/features/quizzes/exams/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamSeed struct {
	NameExams     string  `json:"name_exams"`
	Status        string  `json:"status"`
	TotalQuestion []int64 `json:"total_question"`
	IconURL       string  `json:"icon_url"`
	UnitID        uint    `json:"unit_id"`
	CreatedBy     string  `json:"created_by"`
}

func SeedExamsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []ExamSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.ExamModel
		if err := db.Where("name_exams = ?", seed.NameExams).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Exam '%s' sudah ada, lewati...", seed.NameExams)
			continue
		}

		exam := model.ExamModel{
			NameExams:     seed.NameExams,
			Status:        seed.Status,
			TotalQuestion: seed.TotalQuestion,
			IconURL:       &seed.IconURL,
			UnitID:        seed.UnitID,
			CreatedBy:     parseUUID(seed.CreatedBy),
		}

		if err := db.Create(&exam).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.NameExams, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.NameExams)
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
