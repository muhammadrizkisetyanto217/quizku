package route

import (
	evaluationController "quizku/internals/features/quizzes/evaluations/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EvaluationAdminRoutes(api fiber.Router, db *gorm.DB) {
	evaluationCtrl := evaluationController.NewEvaluationController(db)

	evaluationRoutes := api.Group("/evaluations",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola evaluasi"),
			constants.TeacherAndAbove,
		),
	)
	evaluationRoutes.Post("/", evaluationCtrl.CreateEvaluation)
	evaluationRoutes.Put("/:id", evaluationCtrl.UpdateEvaluation)
	evaluationRoutes.Delete("/:id", evaluationCtrl.DeleteEvaluation)
}
