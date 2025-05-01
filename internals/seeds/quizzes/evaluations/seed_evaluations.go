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
	NameEvaluation string  `json:"name_evaluation"`
	Status         string  `json:"status"`
	TotalQuestion  []int64 `json:"total_question"`
	IconURL        string  `json:"icon_url"`
	UnitID         uint    `json:"unit_id"`
	CreatedBy      string  `json:"created_by"`
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
		if err := db.Where("name_evaluation = ?", seed.NameEvaluation).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Evaluation '%s' sudah ada, lewati...", seed.NameEvaluation)
			continue
		}

		eval := model.EvaluationModel{
			NameEvaluation: seed.NameEvaluation,
			Status:         seed.Status,
			TotalQuestion:  seed.TotalQuestion,
			IconURL:        &seed.IconURL,
			UnitID:         seed.UnitID,
			CreatedBy:      parseUUID(seed.CreatedBy),
		}

		if err := db.Create(&eval).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.NameEvaluation, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.NameEvaluation)
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
