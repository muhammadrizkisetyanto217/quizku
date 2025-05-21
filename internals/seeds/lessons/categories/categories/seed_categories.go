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
	CategoryName             string `json:"category_name"`
	CategoryStatus           string `json:"category_status"`
	CategoryDescriptionShort string `json:"category_description_short"`
	CategoryDescriptionLong  string `json:"category_description_long"`
	CategoryDifficultyID     uint   ` json:"category_difficulty_id"`
}

func SeedCategoriesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file JSON:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []CategorySeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode isi JSON: %v", err)
	}

	for _, c := range inputs {
		var existing categoryModel.CategoryModel
		err := db.Where("category_name = ? AND category_difficulty_id = ?", c.CategoryName, c.CategoryDifficultyID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Data '%s' untuk difficulty_id %d sudah ada, dilewati", c.CategoryName, c.CategoryDifficultyID)
			continue
		}

		newCategory := categoryModel.CategoryModel{
			CategoryName:               c.CategoryName,
			CategoryStatus:             c.CategoryStatus,
			CategoryDescriptionShort:   c.CategoryDescriptionShort,
			CategoryDescriptionLong:    c.CategoryDescriptionLong,
			CategoryDifficultyID:       c.CategoryDifficultyID,
			CategoryTotalSubcategories: pq.Int64Array{},
			CategoryImageURL:           "",
		}

		if err := db.Create(&newCategory).Error; err != nil {
			log.Printf("❌ Gagal insert kategori '%s': %v", c.CategoryName, err)
		} else {
			log.Printf("✅ Berhasil insert kategori '%s'", c.CategoryName)
		}
	}
}
