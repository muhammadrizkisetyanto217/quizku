package unit

import (
	"encoding/json"
	"log"
	"os"
	unitModel "quizku/internals/features/lessons/units/model"

	"gorm.io/gorm"
)

type UnitNewsSeedInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	UnitID      int    `json:"unit_id"`
}

func SeedUnitsNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []UnitNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, news := range inputs {
		var existing unitModel.UnitNewsModel
		err := db.Where("title = ? AND unit_id = ?", news.Title, news.UnitID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ News '%s' untuk unit_id '%d' sudah ada, lewati...", news.Title, news.UnitID)
			continue
		}

		newsEntry := unitModel.UnitNewsModel{
			Title:       news.Title,
			Description: news.Description,
			IsPublic:    news.IsPublic,
			UnitID:      news.UnitID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.Title, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.Title)
		}
	}
}
