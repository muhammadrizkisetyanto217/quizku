package difficulty

import (
	"encoding/json"
	"log"
	"os"
	difficultyModel "quizku/internals/features/lessons/difficulty/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DifficultySeedInput struct {
	DifficultyName             string `json:"difficulty_name"`
	DifficultyStatus           string `json:"difficulty_status"`
	DifficultyDescriptionShort string `json:"difficulty_description_short"`
	DifficultyDescriptionLong  string `json:"difficulty_description_long"`
}

func SeedDifficultiesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []DifficultySeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, input := range inputs {
		var existing difficultyModel.DifficultyModel
		err := db.Where("difficulty_name = ?", input.DifficultyName).First(&existing).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("❌ Gagal cek duplikasi '%s': %v", input.DifficultyName, err)
			continue
		}

		if err == nil {
			log.Printf("ℹ️ Data '%s' sudah ada, lewati...", input.DifficultyName)
			continue
		}

		newEntry := difficultyModel.DifficultyModel{
			DifficultyName:             input.DifficultyName,
			DifficultyStatus:           input.DifficultyStatus,
			DifficultyDescriptionShort: input.DifficultyDescriptionShort,
			DifficultyDescriptionLong:  input.DifficultyDescriptionLong,
			DifficultyTotalCategories:  pq.Int64Array{},
			DifficultyImageURL:         "",
		}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newEntry).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", input.DifficultyName, err)
		} else {
			log.Printf("✅ Berhasil insert difficulty '%s'", input.DifficultyName)
		}
	}
}
