package unit

import (
	"encoding/json"
	"log"
	"os"
	unitModel "quizku/internals/features/lessons/units/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type UnitNewsSeedInput struct {
	UnitNewsTitle       string `json:"title"`
	UnitNewsDescription string `json:"description"`
	UnitNewsIsPublic    bool   `json:"is_public"`
	UnitNewsUnitID      uint   `json:"unit_id"`
}

func SeedUnitsNewsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	// 1. Baca file JSON
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var inputs []UnitNewsSeedInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	// 2. Ambil semua yang sudah ada dari DB
	var existingEntries []unitModel.UnitNewsModel
	if err := db.Select("unit_news_title", "unit_news_unit_id").Find(&existingEntries).Error; err != nil {
		log.Fatalf("❌ Gagal query data existing: %v", err)
	}

	existingMap := make(map[string]bool)
	for _, e := range existingEntries {
		key := e.UnitNewsTitle + "_" + strconv.Itoa(int(e.UnitNewsUnitID))
		existingMap[key] = true
	}

	// 3. Siapkan data untuk bulk insert
	var toInsert []unitModel.UnitNewsModel
	now := time.Now()
	for _, news := range inputs {
		key := news.UnitNewsTitle + "_" + strconv.Itoa(int(news.UnitNewsUnitID))
		if existingMap[key] {
			log.Printf("ℹ️ News '%s' untuk unit_id '%d' sudah ada, lewati...", news.UnitNewsTitle, news.UnitNewsUnitID)
			continue
		}

		toInsert = append(toInsert, unitModel.UnitNewsModel{
			UnitNewsTitle:       news.UnitNewsTitle,
			UnitNewsDescription: news.UnitNewsDescription,
			UnitNewsIsPublic:    news.UnitNewsIsPublic,
			UnitNewsUnitID:      news.UnitNewsUnitID,
			UpdatedAt:           now,
		})
	}

	// 4. Jalankan bulk insert
	if len(toInsert) > 0 {
		if err := db.Create(&toInsert).Error; err != nil {
			log.Fatalf("❌ Gagal bulk insert unit news: %v", err)
		}
		log.Printf("✅ Berhasil insert %d unit news", len(toInsert))
	} else {
		log.Println("ℹ️ Tidak ada unit news baru untuk diinsert.")
	}
}
