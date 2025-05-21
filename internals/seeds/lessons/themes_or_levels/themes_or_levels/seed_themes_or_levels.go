package themes_or_levels

import (
	"encoding/json"
	"log"
	"os"

	themesModel "quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ThemesOrLevelSeedInput struct {
	ThemesOrLevelName             string        `json:"themes_or_level_name"`
	ThemesOrLevelStatus           string        `json:"themes_or_level_status"`
	ThemesOrLevelDescriptionShort string        `json:"themes_or_level_description_short"`
	ThemesOrLevelDescriptionLong  string        `json:"themes_or_level_description_long"`
	ThemesOrLevelTotalUnit        pq.Int64Array `json:"themes_or_level_total_unit"`
	ThemesOrLevelImageURL         string        `json:"themes_or_level_image_url"`
	ThemesOrLevelSubcategoryID    int           `json:"themes_or_level_subcategory_id"`
}

func SeedThemesOrLevelsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	// 1. Baca file
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	// 2. Decode JSON
	var inputs []ThemesOrLevelSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode isi JSON: %v", err)
	}
	log.Printf("📦 Total data themes: %d", len(inputs))

	// 3. Proses setiap data
	for _, input := range inputs {
		if input.ThemesOrLevelName == "" || input.ThemesOrLevelSubcategoryID == 0 {
			log.Printf("⚠️ Dilewati karena data tidak valid: %+v", input)
			continue
		}

		log.Printf("🔍 Cek theme: %s - Subcategory ID: %d", input.ThemesOrLevelName, input.ThemesOrLevelSubcategoryID)

		var existing themesModel.ThemesOrLevelsModel
		err := db.Where("themes_or_level_name = ? AND themes_or_level_subcategory_id = ?", input.ThemesOrLevelName, input.ThemesOrLevelSubcategoryID).
			First(&existing).Error

		if err == nil {
			log.Printf("ℹ️ Theme '%s' sudah ada di subcategory %d, dilewati.", input.ThemesOrLevelName, input.ThemesOrLevelSubcategoryID)
			continue
		}

		newTheme := themesModel.ThemesOrLevelsModel{
			ThemesOrLevelName:             input.ThemesOrLevelName,
			ThemesOrLevelStatus:           input.ThemesOrLevelStatus,
			ThemesOrLevelDescriptionShort: input.ThemesOrLevelDescriptionShort,
			ThemesOrLevelDescriptionLong:  input.ThemesOrLevelDescriptionLong,
			ThemesOrLevelTotalUnit:        input.ThemesOrLevelTotalUnit,
			ThemesOrLevelImageURL:         input.ThemesOrLevelImageURL,
			ThemesOrLevelSubcategoryID:    input.ThemesOrLevelSubcategoryID,
		}

		if err := db.Create(&newTheme).Error; err != nil {
			log.Printf("❌ Gagal insert theme '%s': %v", input.ThemesOrLevelName, err)
		} else {
			log.Printf("✅ Berhasil insert theme '%s'", input.ThemesOrLevelName)
		}
	}
}
