package route

import (
	"quizku/internals/constants"
	"quizku/internals/features/quizzes/questions/controller"
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

	questionLinkController := controller.NewQuestionLinkController(db)

	linkRoutes := api.Group("/question-links",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola relasi soal"),
			constants.TeacherAndAbove,
		),
	)

	linkRoutes.Post("/", questionLinkController.Create)
	linkRoutes.Get("/", questionLinkController.GetAll)
	linkRoutes.Get("/question/:id", questionLinkController.GetByQuestionID)
	linkRoutes.Put("/:id", questionLinkController.Update)
	linkRoutes.Delete("/:id", questionLinkController.Delete)
}
