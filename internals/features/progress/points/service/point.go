package service

import (
	"log"
	levelRequirement "quizku/internals/features/progress/level_rank/model"
	userLogPoint "quizku/internals/features/progress/points/model"
	userProgress "quizku/internals/features/progress/progress/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddUserPointLogAndUpdateProgress(db *gorm.DB, userID uuid.UUID, sourceType int, sourceID int, points int) error {
	log.Printf("[SERVICE] AddUserPointLogAndUpdateProgress - userID: %s sourceType: %d sourceID: %d point: %d",
		userID.String(), sourceType, sourceID, points)

	// 1. Simpan log poin
	logEntry := userLogPoint.UserPointLog{
		UserID:     userID,
		Points:     points,
		SourceType: sourceType,
		SourceID:   sourceID,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log:", err)
		return err
	}

	// 2. Tambahkan poin ke user_progress
	if err := db.Model(&userProgress.UserProgress{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_points": gorm.Expr("total_points + ?", points),
			"last_updated": time.Now(),
		}).Error; err != nil {
		log.Println("[ERROR] Gagal update user_progress:", err)
		return err
	}

	// 3. Ambil user_progress terbaru
	var progress userProgress.UserProgress
	if err := db.Where("user_id = ?", userID).First(&progress).Error; err != nil {
		log.Println("[ERROR] Gagal ambil user_progress setelah update:", err)
		return err
	}

	// 4. Ambil level requirement berdasarkan total_points
	var level levelRequirement.LevelRequirement
	if err := db.Where("min_points <= ? AND (max_points IS NULL OR max_points >= ?)", progress.TotalPoints, progress.TotalPoints).
		Order("level DESC").First(&level).Error; err != nil {
		log.Println("[ERROR] Gagal cari level yang sesuai:", err)
		return err
	}

	// 5. Update level jika naik
	if level.Level != progress.Level {
		if err := db.Model(&userProgress.UserProgress{}).
			Where("user_id = ?", userID).
			Update("level", level.Level).Error; err != nil {
			log.Println("[ERROR] Gagal update level user_progress:", err)
			return err
		}
		log.Printf("[LEVEL-UP] User %s naik ke level %d", userID.String(), level.Level)
		// update struct lokal
		progress.Level = level.Level
	}

	// 6. Ambil rank requirement berdasarkan level sekarang
	var rank levelRequirement.RankRequirement
	if err := db.Where("min_level <= ? AND (max_level IS NULL OR max_level >= ?)", progress.Level, progress.Level).
		Order("rank DESC").First(&rank).Error; err != nil {
		log.Println("[ERROR] Gagal cari rank yang sesuai:", err)
		return err
	}

	// 7. Update rank tanpa if (selalu sinkron dengan level)
	if err := db.Model(&userProgress.UserProgress{}).
		Where("user_id = ?", userID).
		Update("rank", rank.Rank).Error; err != nil {
		log.Println("[ERROR] Gagal update rank user_progress:", err)
		return err
	}
	log.Printf("[RANK-UP] User %s naik ke rank %d (%s)", userID.String(), rank.Rank, rank.Name)

	log.Printf("[SUCCESS] Poin berhasil ditambahkan: %d poin", points)
	return nil
}
