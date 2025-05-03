package route

import (
	quizzesController "quizku/internals/features/quizzes/quizzes/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuizzesAdminRoutes(api fiber.Router, db *gorm.DB) {
	sectionQuizzesController := quizzesController.NewSectionQuizController(db)
	quizController := quizzesController.NewQuizController(db)

	// Section Quizzes
	sectionRoutes := api.Group("/section-quizzes",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola section quiz"),
			constants.TeacherAndAbove,
		),
	)
	sectionRoutes.Post("/", sectionQuizzesController.CreateSectionQuiz)
	sectionRoutes.Put("/:id", sectionQuizzesController.UpdateSectionQuiz)
	sectionRoutes.Delete("/:id", sectionQuizzesController.DeleteSectionQuiz)

	// Quizzes
	quizRoutes := api.Group("/quizzes",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola quiz"),
			constants.TeacherAndAbove,
		),
	)
	quizRoutes.Post("/", quizController.CreateQuiz)
	quizRoutes.Put("/:id", quizController.UpdateQuiz)
	quizRoutes.Delete("/:id", quizController.DeleteQuiz)
}
