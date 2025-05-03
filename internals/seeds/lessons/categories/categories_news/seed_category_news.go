package category

import (
	"encoding/json"
	"log"
	"os"
	categoryModel "quizku/internals/features/lessons/categories/model"

	"gorm.io/gorm"
)

type CategoryNewsSeedInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	CategoryID  int    `json:"category_id"`
}

func SeedCategoriesNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []CategoryNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, news := range inputs {
		var existing categoryModel.CategoryNewsModel
		err := db.Where("title = ? AND category_id = ?", news.Title, news.CategoryID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Data news '%s' untuk category_id '%d' sudah ada, lewati...", news.Title, news.CategoryID)
			continue
		}

		newsEntry := categoryModel.CategoryNewsModel{
			Title:       news.Title,
			Description: news.Description,
			IsPublic:    news.IsPublic,
			CategoryID:  news.CategoryID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.Title, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.Title)
		}
	}
}
