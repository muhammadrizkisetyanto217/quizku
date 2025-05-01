package route

import (
	"quizku/internals/constants"
	evaluationController "quizku/internals/features/quizzes/evaluations/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EvaluationRoute(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🏆 Evaluation Routes
	evaluationCtrl := evaluationController.NewEvaluationController(db)
	evaluationRoutes := api.Group("/evaluations")

	// ✅ GET evaluation → semua user login
	evaluationRoutes.Get("/", evaluationCtrl.GetEvaluations)
	evaluationRoutes.Get("/:id", evaluationCtrl.GetEvaluation)
	evaluationRoutes.Get("/unit/:unitId", evaluationCtrl.GetEvaluationsByUnitID)

	// 🔒 POST, PUT, DELETE → hanya teacher/admin/owner
	evaluationRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("membuat evaluasi"),
		constants.TeacherAndAbove,
	), evaluationCtrl.CreateEvaluation)

	evaluationRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit evaluasi"),
		constants.TeacherAndAbove,
	), evaluationCtrl.UpdateEvaluation)

	evaluationRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus evaluasi"),
		constants.TeacherAndAbove,
	), evaluationCtrl.DeleteEvaluation)

	// 🧠 User Evaluation Routes → semua user login
	userEvaluationController := evaluationController.NewUserEvaluationController(db)
	userEvaluationRoutes := api.Group("/user-evaluations")
	userEvaluationRoutes.Post("/", userEvaluationController.Create)
	userEvaluationRoutes.Get("/:user_id", userEvaluationController.GetByUserID)
}
