package route

import (
	readingController "quizku/internals/features/quizzes/readings/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReadingAdminRoutes(api fiber.Router, db *gorm.DB) {
	readingCtrl := readingController.NewReadingController(db)

	readingRoutes := api.Group("/readings",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola reading"),
			constants.TeacherAndAbove,
		),
	)
	readingRoutes.Post("/", readingCtrl.CreateReading)
	readingRoutes.Put("/:id", readingCtrl.UpdateReading)
	readingRoutes.Delete("/:id", readingCtrl.DeleteReading)
}
