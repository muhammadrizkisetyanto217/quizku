package reading

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/quizzes/readings/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingSeed struct {
	ReadingTitle           string `json:"reading_title"`
	ReadingStatus          string `json:"reading_status"`
	ReadingDescriptionLong string `json:"reading_description_long"`
	ReadingUnitID          uint   `json:"reading_unit_id"`
	ReadingCreatedBy       string `json:"reading_created_by"`
}

func SeedReadingsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []ReadingSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.ReadingModel
		if err := db.Where("reading_title = ?", seed.ReadingTitle).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Reading '%s' sudah ada, lewati...", seed.ReadingTitle)
			continue
		}

		reading := model.ReadingModel{
			ReadingTitle:           seed.ReadingTitle,
			ReadingStatus:          seed.ReadingStatus,
			ReadingDescriptionLong: seed.ReadingDescriptionLong,
			ReadingUnitID:          seed.ReadingUnitID,
			ReadingCreatedBy:       parseUUID(seed.ReadingCreatedBy),
		}

		if err := db.Create(&reading).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.ReadingTitle, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.ReadingTitle)
		}
	}
}

func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("❌ Gagal parse UUID: %v", err)
	}
	return id
}
