package levels

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/progress/level_rank/model"

	"gorm.io/gorm"
)

type LevelSeed struct {
	Level     int    `json:"level"`
	Name      string `json:"name"`
	MinPoints int    `json:"min_points"`
	MaxPoints *int   `json:"max_points"`
}

func SeedLevelRequirementsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal baca file JSON: %v", err)
	}

	var data []LevelSeed
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, item := range data {
		var existing model.LevelRequirement
		if err := db.Where("level = ?", item.Level).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Level %d sudah ada, lewati...", item.Level)
			continue
		}

		record := model.LevelRequirement{
			Level:     item.Level,
			Name:      item.Name,
			MinPoints: item.MinPoints,
			MaxPoints: item.MaxPoints,
		}

		if err := db.Create(&record).Error; err != nil {
			log.Printf("❌ Gagal insert Level %d: %v", item.Level, err)
		} else {
			log.Printf("✅ Berhasil insert Level %d", item.Level)
		}
	}
}
