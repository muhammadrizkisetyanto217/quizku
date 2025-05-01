package rank

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/progress/level_rank/model"

	"gorm.io/gorm"
)

type RankSeed struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	MinLevel int    `json:"min_level"`
	MaxLevel *int   `json:"max_level"` // nullable
}

func SeedRanksRequirementsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var input []RankSeed
	if err := json.Unmarshal(file, &input); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, r := range input {
		var existing model.RankRequirement
		if err := db.Where("rank = ?", r.Rank).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Rank %d sudah ada, lewati...", r.Rank)
			continue
		}

		newRank := model.RankRequirement{
			Rank:     r.Rank,
			Name:     r.Name,
			MinLevel: r.MinLevel,
			MaxLevel: r.MaxLevel,
		}

		if err := db.Create(&newRank).Error; err != nil {
			log.Printf("❌ Gagal insert Rank %d: %v", r.Rank, err)
		} else {
			log.Printf("✅ Berhasil insert Rank %d (%s)", r.Rank, r.Name)
		}
	}
}
