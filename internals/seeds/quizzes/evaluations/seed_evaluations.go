package evaluation

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/quizzes/evaluations/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EvaluationSeed struct {
	EvaluationName          string  `json:"evaluation_name"`
	EvaluationStatus        string  `json:"evaluation_status"`
	EvaluationTotalQuestion []int64 `json:"evaluation_total_question"`
	EvaluationIconURL       string  `json:"evaluation_icon_url"`
	EvaluationUnitID        uint    `json:"evaluation_unit_id"`
	EvaluationCreatedBy     string  `json:"evaluation_created_by"` // UUID dalam bentuk string
}

func SeedEvaluationsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []EvaluationSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.EvaluationModel
		if err := db.Where("evaluation_name = ?", seed.EvaluationName).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Evaluation '%s' sudah ada, lewati...", seed.EvaluationName)
			continue
		}

		eval := model.EvaluationModel{
			EvaluationName:          seed.EvaluationName,
			EvaluationStatus:        seed.EvaluationStatus,
			EvaluationTotalQuestion: seed.EvaluationTotalQuestion,
			EvaluationIconURL:       &seed.EvaluationIconURL,
			EvaluationUnitID:        seed.EvaluationUnitID,
			EvaluationCreatedBy:     parseUUID(seed.EvaluationCreatedBy),
		}

		if err := db.Create(&eval).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.EvaluationName, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.EvaluationName)
		}
	}
}

// Helper UUID validasi
func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("❌ UUID tidak valid: %v", err)
	}
	return id
}
