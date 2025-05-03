package subcategory

import (
	"encoding/json"
	"log"
	"os"
	subcategoryModel "quizku/internals/features/lessons/subcategory/model"

	"gorm.io/gorm"
)

type SubcategoryNewsSeedInput struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	IsPublic       bool   `json:"is_public"`
	SubcategoriesID uint   `json:"subcategories_id"`
}

func SeedSubcategoriesNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []SubcategoryNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, news := range inputs {
		var existing subcategoryModel.SubcategoryNewsModel
		err := db.Where("title = ? AND subcategory_id = ?", news.Title, news.SubcategoriesID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Data news '%s' untuk subcategory_id '%d' sudah ada, lewati...", news.Title, news.SubcategoriesID)
			continue
		}

		newsEntry := subcategoryModel.SubcategoryNewsModel{
			Title:           news.Title,
			Description:     news.Description,
			IsPublic:        news.IsPublic,
			SubcategoryID: news.SubcategoriesID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.Title, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.Title)
		}
	}
}
