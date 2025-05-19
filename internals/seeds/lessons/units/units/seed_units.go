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
	UnitName                string    `json:"unit_name"`
	UnitStatus              string    `json:"unit_status"`
	UnitDescriptionShort    string    `json:"unit_description_short"`
	UnitDescriptionOverview string    `json:"unit_description_overview"`
	UnitImageURL            string    `json:"unit_image_url"`
	UnitTotalSectionQuizzes []int64   `json:"unit_total_section_quizzes"`
	UnitThemesOrLevelID     uint      `json:"unit_themes_or_level_id"`
	UnitCreatedBy           uuid.UUID `json:"unit_created_by"`
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
		if err := db.Where("unit_name = ?", data.UnitName).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Data unit '%s' sudah ada, dilewati.", data.UnitName)
			continue
		}

		newUnit := model.UnitModel{
			UnitName:                data.UnitName,
			UnitStatus:              data.UnitStatus,
			UnitDescriptionShort:    data.UnitDescriptionShort,
			UnitDescriptionOverview: data.UnitDescriptionOverview,
			UnitImageURL:            data.UnitImageURL,
			UnitTotalSectionQuizzes: data.UnitTotalSectionQuizzes,
			UnitThemesOrLevelID:     data.UnitThemesOrLevelID,
			UnitCreatedBy:           data.UnitCreatedBy,
		}

		if err := db.Create(&newUnit).Error; err != nil {
			log.Printf("❌ Gagal insert unit '%s': %v", data.UnitName, err)
		} else {
			log.Printf("✅ Berhasil insert unit '%s'", data.UnitName)
		}
	}
}
