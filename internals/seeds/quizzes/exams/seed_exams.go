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
	ExamName             string  `json:"exam_name"`
	ExamStatus           string  `json:"exam_status"`
	ExamTotalQuestionIDs []int64 `json:"exam_total_question_ids"`
	ExamIconURL          string  `json:"exam_icon_url"`
	ExamUnitID           uint    `json:"exam_unit_id"`
	ExamCreatedBy        string  `json:"exam_created_by"`
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
		if err := db.Where("exam_name = ?", seed.ExamName).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Exam '%s' sudah ada, lewati...", seed.ExamName)
			continue
		}

		exam := model.ExamModel{
			ExamName:             seed.ExamName,
			ExamStatus:           seed.ExamStatus,
			ExamTotalQuestionIDs: seed.ExamTotalQuestionIDs,
			ExamIconURL:          &seed.ExamIconURL,
			ExamUnitID:           seed.ExamUnitID,
			ExamCreatedBy:        parseUUID(seed.ExamCreatedBy),
		}

		if err := db.Create(&exam).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.ExamName, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.ExamName)
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
