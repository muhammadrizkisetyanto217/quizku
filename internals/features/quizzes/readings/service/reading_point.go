package service

import (
	addUserPointService "quizku/internals/features/progress/points/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPointFromReading(db *gorm.DB, userID uuid.UUID, readingID uint, attempt int) error {
	var point int
	switch attempt {
	case 1:
		point = 10
	case 2:
		point = 20
	default:
		point = 5
	}

	const sourceTypeReading = 0

	return addUserPointService.AddUserPointLogAndUpdateProgress(db, userID, sourceTypeReading, int(readingID), point)
}
