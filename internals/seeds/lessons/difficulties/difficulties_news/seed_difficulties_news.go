package difficulty

import (
	"encoding/json"
	"log"
	"os"

	difficultyModel "quizku/internals/features/lessons/difficulty/model"

	"gorm.io/gorm"
)

type DifficultyNewsSeedInput struct {
	DifficultyNewsTitle        string `json:"difficulty_news_title"`
	DifficultyNewsDescription  string `json:"difficulty_news_description"`
	DifficultyNewsIsPublic     bool   `json:"difficulty_news_is_public"`
	DifficultyNewsDifficultyID uint   `json:"difficulty_news_difficulty_id"`
}

func SeedDifficultyNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file JSON:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []DifficultyNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode isi JSON: %v", err)
	}

	for _, input := range inputs {
		var existing difficultyModel.DifficultyNewsModel
		err := db.Where(
			"difficulty_news_title = ? AND difficulty_news_difficulty_id = ?",
			input.DifficultyNewsTitle, input.DifficultyNewsDifficultyID,
		).First(&existing).Error

		if err == nil {
			log.Printf("ℹ️ Data '%s' untuk difficulty_id %d sudah ada, dilewati...", input.DifficultyNewsTitle, input.DifficultyNewsDifficultyID)
			continue
		}

		newsEntry := difficultyModel.DifficultyNewsModel{
			DifficultyNewsTitle:        input.DifficultyNewsTitle,
			DifficultyNewsDescription:  input.DifficultyNewsDescription,
			DifficultyNewsIsPublic:     input.DifficultyNewsIsPublic,
			DifficultyNewsDifficultyID: input.DifficultyNewsDifficultyID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", input.DifficultyNewsTitle, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", input.DifficultyNewsTitle)
		}
	}
}
