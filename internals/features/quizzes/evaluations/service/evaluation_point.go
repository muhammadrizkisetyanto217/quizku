package service

import (
	addUserPointService "quizku/internals/features/progress/points/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPointFromEvaluation(db *gorm.DB, userID uuid.UUID, evaluationID uint, attempt int) error {
	var point int
	switch attempt {
	case 1:
		point = 25
	case 2:
		point = 15
	default:
		point = 10
	}

	const sourceTypeEvaluation = 2

	return addUserPointService.AddUserPointLogAndUpdateProgress(db, userID, sourceTypeEvaluation, int(evaluationID), point)
}
