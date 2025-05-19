package themes

import (
	"encoding/json"
	"log"
	"os"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"

	"gorm.io/gorm"
)

// ✅ Struct input sesuai penamaan semantik di model
type ThemesOrLevelsNewsSeedInput struct {
	ThemesNewsTitle           string `json:"themes_news_title"`
	ThemesNewsDescription     string `json:"themes_news_description"`
	ThemesNewsIsPublic        bool   `json:"themes_news_is_public"`
	ThemesNewsThemesOrLevelID uint   `json:"themes_news_themes_or_level_id"`
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

	for _, input := range inputs {
		var existing themesModel.ThemesOrLevelsNewsModel
		err := db.Where("themes_news_title = ? AND themes_news_themes_or_level_id = ?", input.ThemesNewsTitle, input.ThemesNewsThemesOrLevelID).
			First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ News '%s' untuk themes_or_level_id '%d' sudah ada, dilewati...", input.ThemesNewsTitle, input.ThemesNewsThemesOrLevelID)
			continue
		}

		newsEntry := themesModel.ThemesOrLevelsNewsModel{
			ThemesNewsTitle:           input.ThemesNewsTitle,
			ThemesNewsDescription:     input.ThemesNewsDescription,
			ThemesNewsIsPublic:        input.ThemesNewsIsPublic,
			ThemesNewsThemesOrLevelID: input.ThemesNewsThemesOrLevelID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", input.ThemesNewsTitle, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", input.ThemesNewsTitle)
		}
	}
}
