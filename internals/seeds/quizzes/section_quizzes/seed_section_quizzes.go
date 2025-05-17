package sectionquizzes

import (
	"encoding/json"
	"log"
	"os"

	sectionQuizModel "quizku/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SectionQuizSeed struct {
	SectionQuizzesName         string  `json:"section_quizzes_name"`
	SectionQuizzesStatus       string  `json:"section_quizzes_status"`
	SectionQuizzesMaterials    string  `json:"section_quizzes_materials"`
	SectionQuizzesIconURL      string  `json:"section_quizzes_icon_url"`
	SectionQuizzesUnitID       uint    `json:"section_quizzes_unit_id"`
	SectionQuizzesCreatedBy    string  `json:"section_quizzes_created_by"`
	SectionQuizzesTotalQuizzes []int64 `json:"section_quizzes_total_quizzes"`
}

func SeedSectionQuizzesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []SectionQuizSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing sectionQuizModel.SectionQuizzesModel
		if err := db.Where("section_quizzes_name = ?", seed.SectionQuizzesName).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Section Quiz '%s' sudah ada, lewati...", seed.SectionQuizzesName)
			continue
		}

		newSection := sectionQuizModel.SectionQuizzesModel{
			SectionQuizzesName:         seed.SectionQuizzesName,
			SectionQuizzesStatus:       seed.SectionQuizzesStatus,
			SectionQuizzesMaterials:    seed.SectionQuizzesMaterials,
			SectionQuizzesIconURL:      seed.SectionQuizzesIconURL,
			SectionQuizzesTotalQuizzes: pq.Int64Array(seed.SectionQuizzesTotalQuizzes),
			SectionQuizzesUnitID:       seed.SectionQuizzesUnitID,
			SectionQuizzesCreatedBy:    parseUUID(seed.SectionQuizzesCreatedBy),
		}

		if err := db.Create(&newSection).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.SectionQuizzesName, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.SectionQuizzesName)
		}
	}
}

// helper mandiri
func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("❌ Gagal parse UUID: %v", err)
	}
	return id
}
