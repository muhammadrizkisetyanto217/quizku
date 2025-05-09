package tooltip

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/utils/tooltips/model"

	"gorm.io/gorm"
)

type TooltipSeed struct {
	Keyword          string `json:"keyword"`
	DescriptionShort string `json:"description_short"`
	DescriptionLong  string `json:"description_long"`
}

func SeedTooltipsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []TooltipSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.Tooltip
		if err := db.Where("keyword = ?", seed.Keyword).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Tooltip '%s' sudah ada, lewati...", seed.Keyword)
			continue
		}

		tooltip := model.Tooltip{
			Keyword:          seed.Keyword,
			DescriptionShort: seed.DescriptionShort,
			DescriptionLong:  seed.DescriptionLong,
		}

		if err := db.Create(&tooltip).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.Keyword, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.Keyword)
		}
	}
}
