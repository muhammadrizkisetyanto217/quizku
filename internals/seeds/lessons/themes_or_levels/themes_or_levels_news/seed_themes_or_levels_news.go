package themes

import (
	"encoding/json"
	"log"
	"os"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"

	"gorm.io/gorm"
)

type ThemesOrLevelsNewsSeedInput struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	IsPublic         bool   `json:"is_public"`
	ThemesOrLevelsID uint   `json:"themes_or_levels_id"`
}

func SeedThemesOrLevelsNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []ThemesOrLevelsNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, news := range inputs {
		var existing themesModel.ThemesOrLevelsNewsModel
		err := db.Where("title = ? AND themes_or_levels_id = ?", news.Title, news.ThemesOrLevelsID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ News '%s' untuk themes_or_levels_id '%d' sudah ada, lewati...", news.Title, news.ThemesOrLevelsID)
			continue
		}

		newsEntry := themesModel.ThemesOrLevelsNewsModel{
			Title:            news.Title,
			Description:      news.Description,
			IsPublic:         news.IsPublic,
			ThemesOrLevelsID: news.ThemesOrLevelsID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.Title, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.Title)
		}
	}
}
