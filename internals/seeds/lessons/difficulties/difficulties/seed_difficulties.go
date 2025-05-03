package difficulty

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/lessons/difficulty/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DifficultySeed struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	DescriptionShort string `json:"description_short"`
	DescriptionLong  string `json:"description_long"`
}

func SeedDifficultiesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var input []DifficultySeed
	if err := json.Unmarshal(file, &input); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, d := range input {
		var existing model.DifficultyModel
		err := db.Where("name = ?", d.Name).First(&existing).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("❌ Error saat cek data '%s': %v", d.Name, err)
			continue
		}

		if err == nil {
			log.Printf("ℹ️ Data dengan nama '%s' sudah ada, lewati...", d.Name)
			continue
		}

		newDifficulty := model.DifficultyModel{
			Name:             d.Name,
			Status:           d.Status,
			DescriptionShort: d.DescriptionShort,
			DescriptionLong:  d.DescriptionLong,
		}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newDifficulty).Error; err != nil {
			log.Printf("❌ Gagal insert data '%s': %v", d.Name, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", d.Name)
		}
	}
}
