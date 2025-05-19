package service

import (
	addUserPointService "quizku/internals/features/progress/points/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPointFromExam(db *gorm.DB, userID uuid.UUID, examID uint, attempt int) error {
	var point int
	switch attempt {
	case 1:
		point = 20
	case 2:
		point = 40
	default:
		point = 10
	}
	return addUserPointService.AddUserPointLogAndUpdateProgress(db, userID, 3, int(examID), point)
}