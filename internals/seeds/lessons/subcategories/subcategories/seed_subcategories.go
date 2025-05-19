package subcategory

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/lessons/subcategories/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SubcategorySeed struct {
	Name                string  `json:"name"`
	Status              string  `json:"status"`
	DescriptionLong     string  `json:"description_long"`
	ImageURL            string  `json:"image_url"`
	TotalThemesOrLevels []int64 `json:"total_themes_or_levels"`
	CategoriesID        uint    `json:"categories_id"`
}

func SeedSubcategoriesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var input []SubcategorySeed
	if err := json.Unmarshal(file, &input); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, s := range input {
		var existing model.SubcategoryModel
		if err := db.Where("subcategory_name = ?", s.Name).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Subkategori '%s' sudah ada, lewati...", s.Name)
			continue
		}

		sub := model.SubcategoryModel{
			SubcategoryName:                s.Name,
			SubcategoryStatus:              s.Status,
			SubcategoryDescriptionLong:     s.DescriptionLong,
			SubcategoryImageURL:            s.ImageURL,
			SubcategoryTotalThemesOrLevels: pq.Int64Array(s.TotalThemesOrLevels),
			SubcategoryCategoryID:          s.CategoriesID,
		}

		if err := db.Create(&sub).Error; err != nil {
			log.Printf("❌ Gagal insert subkategori '%s': %v", s.Name, err)
		} else {
			log.Printf("✅ Berhasil insert subkategori '%s'", s.Name)
		}
	}
}
