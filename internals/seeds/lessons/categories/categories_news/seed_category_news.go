package category

import (
	"encoding/json"
	"log"
	"os"

	categoryModel "quizku/internals/features/lessons/categories/model"

	"gorm.io/gorm"
)

type CategoryNewsSeedInput struct {
	CategoryNewsTitle       string `json:"category_news_title"`
	CategoryNewsDescription string `json:"category_news_description"`
	CategoryNewsIsPublic    bool   `json:"category_news_is_public"`
	CategoryNewsCategoryID  uint   `json:"category_news_category_id"`
}

func SeedCategoriesNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file JSON:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file: %v", err)
	}

	var inputs []CategoryNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal parsing JSON: %v", err)
	}

	for _, input := range inputs {
		var existing categoryModel.CategoryNewsModel
		err := db.Where("category_news_title = ? AND category_news_category_id = ?", input.CategoryNewsTitle, input.CategoryNewsCategoryID).
			First(&existing).Error

		if err == nil {
			log.Printf("ℹ️ Data '%s' untuk category_id %d sudah ada, dilewati", input.CategoryNewsTitle, input.CategoryNewsCategoryID)
			continue
		}

		news := categoryModel.CategoryNewsModel{
			CategoryNewsTitle:       input.CategoryNewsTitle,
			CategoryNewsDescription: input.CategoryNewsDescription,
			CategoryNewsIsPublic:    input.CategoryNewsIsPublic,
			CategoryNewsCategoryID:  input.CategoryNewsCategoryID,
		}

		if err := db.Create(&news).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", input.CategoryNewsTitle, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", input.CategoryNewsTitle)
		}
	}
}
