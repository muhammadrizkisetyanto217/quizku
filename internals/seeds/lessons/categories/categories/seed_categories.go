package category

import (
	"encoding/json"
	"log"
	"os"
	categoryModel "quizku/internals/features/lessons/categories/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CategorySeedInput struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	DescriptionShort string `json:"description_short"`
	DescriptionLong  string `json:"description_long"`
	DifficultyID     uint   `json:"difficulty_id"`
}

func SeedCategoriesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []CategorySeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, c := range inputs {
		var existing categoryModel.CategoryModel
		err := db.Where("name = ? AND difficulty_id = ?", c.Name, c.DifficultyID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Data dengan nama '%s' dan difficulty_id '%d' sudah ada, lewati...", c.Name, c.DifficultyID)
			continue
		}

		newCategory := categoryModel.CategoryModel{
			Name:               c.Name,
			Status:             c.Status,
			DescriptionShort:   c.DescriptionShort,
			DescriptionLong:    c.DescriptionLong,
			DifficultyID:       c.DifficultyID,
			TotalSubcategories: pq.Int64Array{},
			ImageURL:           "",
		}

		if err := db.Create(&newCategory).Error; err != nil {
			log.Printf("❌ Gagal insert data '%s': %v", c.Name, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", c.Name)
		}
	}
}
