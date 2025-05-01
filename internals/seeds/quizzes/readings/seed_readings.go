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
	Title           string  `json:"title"`
	Status          string  `json:"status"`
	DescriptionLong string  `json:"description_long"`
	TooltipsID      []int64 `json:"tooltips_id"`
	UnitID          uint    `json:"unit_id"`
	CreatedBy       string  `json:"created_by"`
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
		if err := db.Where("title = ?", seed.Title).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Reading '%s' sudah ada, lewati...", seed.Title)
			continue
		}

		reading := model.ReadingModel{
			Title:           seed.Title,
			Status:          seed.Status,
			DescriptionLong: seed.DescriptionLong,
			TooltipsID:      seed.TooltipsID,
			UnitID:          seed.UnitID,
			CreatedBy:       parseUUID(seed.CreatedBy),
		}

		if err := db.Create(&reading).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.Title, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.Title)
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
