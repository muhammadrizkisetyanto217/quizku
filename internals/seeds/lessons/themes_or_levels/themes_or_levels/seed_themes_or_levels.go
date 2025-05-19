package themes_or_levels

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/lessons/themes_or_levels/model"

	"gorm.io/gorm"
)

type ThemesSeed struct {
	ThemesOrLevelName             string `json:"name"`
	ThemesOrLevelStatus           string `json:"status"`
	ThemesOrLevelDescriptionShort string `json:"description_short"`
	ThemesOrLevelDescriptionLong  string `json:"description_long"`
	ThemesOrLevelImageURL         string `json:"image_url"`
	ThemesOrLevelSubcategoryID    int    `json:"subcategories_id"`
}

func SeedThemesOrLevelsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var input []ThemesSeed
	if err := json.Unmarshal(file, &input); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, t := range input {
		var existing model.ThemesOrLevelsModel
		if err := db.Where("themes_or_level_name = ? AND themes_or_level_subcategory_id = ?", t.ThemesOrLevelName, t.ThemesOrLevelSubcategoryID).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Data '%s' sudah ada untuk subcategory_id %d, dilewati.", t.ThemesOrLevelName, t.ThemesOrLevelSubcategoryID)
			continue
		}

		newTheme := model.ThemesOrLevelsModel{
			ThemesOrLevelName:             t.ThemesOrLevelName,
			ThemesOrLevelStatus:           t.ThemesOrLevelStatus,
			ThemesOrLevelDescriptionShort: t.ThemesOrLevelDescriptionShort,
			ThemesOrLevelDescriptionLong:  t.ThemesOrLevelDescriptionLong,
			ThemesOrLevelSubcategoryID:    t.ThemesOrLevelSubcategoryID,
			ThemesOrLevelImageURL:         t.ThemesOrLevelImageURL,
		}

		if err := db.Create(&newTheme).Error; err != nil {
			log.Printf("❌ Gagal insert theme '%s': %v", t.ThemesOrLevelName, err)
		} else {
			log.Printf("✅ Berhasil insert theme '%s'", t.ThemesOrLevelName)
		}
	}
}
