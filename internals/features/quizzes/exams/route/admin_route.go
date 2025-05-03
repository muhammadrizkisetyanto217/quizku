package route

import (
	examController "quizku/internals/features/quizzes/exams/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExamAdminRoutes(api fiber.Router, db *gorm.DB) {
	examCtrl := examController.NewExamController(db)

	examRoutes := api.Group("/exams",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola ujian"),
			constants.TeacherAndAbove,
		),
	)
	examRoutes.Post("/", examCtrl.CreateExam)
	examRoutes.Put("/:id", examCtrl.UpdateExam)
	examRoutes.Delete("/:id", examCtrl.DeleteExam)
}
