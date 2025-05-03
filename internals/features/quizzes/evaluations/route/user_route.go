package route

import (
	evaluationController "quizku/internals/features/quizzes/evaluations/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EvaluationUserRoutes(api fiber.Router, db *gorm.DB) {
	evaluationCtrl := evaluationController.NewEvaluationController(db)
	userEvaluationCtrl := evaluationController.NewUserEvaluationController(db)

	evaluationRoutes := api.Group("/evaluations")
	evaluationRoutes.Get("/", evaluationCtrl.GetEvaluations)
	evaluationRoutes.Get("/:id", evaluationCtrl.GetEvaluation)
	evaluationRoutes.Get("/unit/:unitId", evaluationCtrl.GetEvaluationsByUnitID)

	userEvaluationRoutes := api.Group("/user-evaluations")
	userEvaluationRoutes.Post("/", userEvaluationCtrl.Create)
	userEvaluationRoutes.Get("/:user_id", userEvaluationCtrl.GetByUserID)
}
