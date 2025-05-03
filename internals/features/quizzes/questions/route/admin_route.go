package route

import (
	"quizku/internals/features/quizzes/questions/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuestionAdminRoutes(api fiber.Router, db *gorm.DB) {
	questionController := controller.NewQuestionController(db)

	questionRoutes := api.Group("/question",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola soal"),
			constants.TeacherAndAbove,
		),
	)
	questionRoutes.Post("/", questionController.CreateQuestion)
	questionRoutes.Put("/:id", questionController.UpdateQuestion)
	questionRoutes.Delete("/:id", questionController.DeleteQuestion)
}
