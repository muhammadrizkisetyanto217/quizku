package subcategory

import (
	"encoding/json"
	"log"
	"os"

	subcategoryModel "quizku/internals/features/lessons/subcategories/model"

	"gorm.io/gorm"
)

type SubcategoryNewsSeedInput struct {
	SubcategoryNewsTitle         string `json:"subcategory_news_title"`
	SubcategoryNewsDescription   string `json:"subcategory_news_description"`
	SubcategoryNewsIsPublic      bool   `json:"subcategory_news_is_public"`
	SubcategoryNewsSubcategoryID uint   `json:"subcategory_news_subcategory_id"`
}

func SeedSubcategoryNewsFromJSON(db *gorm.DB, filePath string) {
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
		err := db.
			Where("subcategory_news_title = ? AND subcategory_news_subcategory_id = ?", news.SubcategoryNewsTitle, news.SubcategoryNewsSubcategoryID).
			First(&existing).Error

		if err == nil {
			log.Printf("ℹ️ Data news '%s' untuk subcategory_id '%d' sudah ada, lewati...",
				news.SubcategoryNewsTitle, news.SubcategoryNewsSubcategoryID)
			continue
		}

		newsEntry := subcategoryModel.SubcategoryNewsModel{
			SubcategoryNewsTitle:         news.SubcategoryNewsTitle,
			SubcategoryNewsDescription:   news.SubcategoryNewsDescription,
			SubcategoryNewsIsPublic:      news.SubcategoryNewsIsPublic,
			SubcategoryNewsSubcategoryID: news.SubcategoryNewsSubcategoryID,
		}

		if err := db.Create(&newsEntry).Error; err != nil {
			log.Printf("❌ Gagal insert news '%s': %v", news.SubcategoryNewsTitle, err)
		} else {
			log.Printf("✅ Berhasil insert news '%s'", news.SubcategoryNewsTitle)
		}
	}
}
