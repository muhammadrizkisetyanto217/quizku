package units

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/lessons/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnitSeed struct {
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	DescriptionShort    string    `json:"description_short"`
	DescriptionOverview string    `json:"description_overview"`
	ImageURL            string    `json:"image_url"`
	TotalSectionQuizzes []int64   `json:"total_section_quizzes"`
	ThemesOrLevelID     uint      `json:"themes_or_level_id"`
	CreatedBy           uuid.UUID `json:"created_by"`
}

func SeedUnitsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file unit:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []UnitSeed
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, data := range inputs {
		// Cek apakah sudah ada berdasarkan nama
		var existing model.UnitModel
		if err := db.Where("name = ?", data.Name).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Data unit '%s' sudah ada, dilewati.", data.Name)
			continue
		}

		newUnit := model.UnitModel{
			Name:                data.Name,
			Status:              data.Status,
			DescriptionShort:    data.DescriptionShort,
			DescriptionOverview: data.DescriptionOverview,
			ImageURL:            data.ImageURL,
			TotalSectionQuizzes: data.TotalSectionQuizzes,
			ThemesOrLevelID:     data.ThemesOrLevelID,
			CreatedBy:           data.CreatedBy,
		}

		if err := db.Create(&newUnit).Error; err != nil {
			log.Printf("❌ Gagal insert unit '%s': %v", data.Name, err)
		} else {
			log.Printf("✅ Berhasil insert unit '%s'", data.Name)
		}
	}
}
