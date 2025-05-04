package scheduler

import (
	"log"
	"os"
	"strconv"
	"time"

	"quizku/internals/features/users/auth/models"

	"gorm.io/gorm"
)

func StartBlacklistCleanupScheduler(db *gorm.DB) {
	go func() {
		// Ambil TTL dari env, default 7 hari
		ttlDays := 7
		if val := os.Getenv("TOKEN_BLACKLIST_TTL_DAYS"); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil {
				ttlDays = parsed
			}
		}

		for {
			log.Println("[CLEANUP] Menjalankan pembersihan token_blacklist...")

			// Hanya hapus token yang expired > ttlDays yang lalu
			deleteBefore := time.Now().Add(-time.Duration(ttlDays) * 24 * time.Hour)
			result := db.Where("expired_at < ?", deleteBefore).Delete(&models.TokenBlacklist{})

			if result.Error != nil {
				log.Printf("[CLEANUP ERROR] %v", result.Error)
			} else if result.RowsAffected > 0 {
				log.Printf("[CLEANUP] %d token kadaluarsa dihapus", result.RowsAffected)
			} else {
				log.Println("[CLEANUP] Tidak ada token yang memenuhi syarat dihapus")
			}

			time.Sleep(24 * time.Hour)
		}
	}()
}
