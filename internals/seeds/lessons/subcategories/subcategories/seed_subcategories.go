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
	SubcategoryName                string  `json:"subcategory_name"`
	SubcategoryStatus              string  `json:"subcategory_status"`
	SubcategoryDescriptionLong     string  `json:"subcategory_description_long"`
	SubcategoryImageURL            string  `json:"subcategory_image_url"`
	SubcategoryTotalThemesOrLevels []int64 `json:"subcategory_total_themes_or_levels"`
	SubcategoryCategoryID          uint    `json:"subcategory_category_id"`
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
	// Tambahan validasi biar aman
	if s.SubcategoryName == "" || s.SubcategoryCategoryID == 0 {
		log.Printf("⚠️ Subkategori kosong atau category_id = 0 dilewati.")
		continue
	}

	var existing model.SubcategoryModel
	if err := db.Where("subcategory_name = ?", s.SubcategoryName).First(&existing).Error; err == nil {
		log.Printf("ℹ️ Subkategori '%s' sudah ada, lewati...", s.SubcategoryName)
		continue
	}

	sub := model.SubcategoryModel{
		SubcategoryName:                s.SubcategoryName,
		SubcategoryStatus:              s.SubcategoryStatus,
		SubcategoryDescriptionLong:     s.SubcategoryDescriptionLong,
		SubcategoryImageURL:            s.SubcategoryImageURL,
		SubcategoryTotalThemesOrLevels: pq.Int64Array(s.SubcategoryTotalThemesOrLevels),
		SubcategoryCategoryID:          s.SubcategoryCategoryID,
	}

	if err := db.Create(&sub).Error; err != nil {
		log.Printf("❌ Gagal insert subkategori '%s': %v", s.SubcategoryName, err)
	} else {
		log.Printf("✅ Berhasil insert subkategori '%s'", s.SubcategoryName)
	}
}

}
