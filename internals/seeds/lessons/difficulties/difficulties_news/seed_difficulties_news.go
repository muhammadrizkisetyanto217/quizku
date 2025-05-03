package difficulty

import (
	"encoding/json"
	"log"
	"os"
	difficultyModel "quizku/internals/features/lessons/difficulty/model"

	"gorm.io/gorm"
)

type DifficultyNewsSeedInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	IsPublic     bool   `json:"is_public"`
	DifficultyID uint   `json:"difficulty_id"`
}

func SeedDifficultiesNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []DifficultyNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, news := range inputs {
		var existing difficultyModel.DifficultyNews
		err := db.Where("title = ? AND difficulty_id = ?", news.Title, news.DifficultyID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Data news '%s' untuk difficulty_id '%d' sudah ada, lewati...", news.Title, news.DifficultyID)
			continue
		}

		newsEntry := difficultyModel.DifficultyNews{
			Title:        news.Title,
			Description:  news.Description,
			IsPublic:     news.IsPublic,
			DifficultyID: news.DifficultyID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.Title, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.Title)
		}
	}
}
