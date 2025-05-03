package themes_or_levels

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/lessons/themes_or_levels/model"

	"gorm.io/gorm"
)

type ThemesSeed struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	DescriptionShort string `json:"description_short"`
	DescriptionLong  string `json:"description_long"`
	ImageURL         string `json:"image_url"`
	SubcategoriesID  int    `json:"subcategories_id"`
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
		if err := db.Where("name = ? AND subcategories_id = ?", t.Name, t.SubcategoriesID).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Data '%s' sudah ada untuk subcategory_id %d, dilewati.", t.Name, t.SubcategoriesID)
			continue
		}

		newTheme := model.ThemesOrLevelsModel{
			Name:             t.Name,
			Status:           t.Status,
			DescriptionShort: t.DescriptionShort,
			DescriptionLong:  t.DescriptionLong,
			SubcategoriesID:  t.SubcategoriesID,
			ImageURL:         t.ImageURL,
		}

		if err := db.Create(&newTheme).Error; err != nil {
			log.Printf("❌ Gagal insert theme '%s': %v", t.Name, err)
		} else {
			log.Printf("✅ Berhasil insert theme '%s'", t.Name)
		}
	}
}
