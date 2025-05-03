package route

import (
	unitController "quizku/internals/features/lessons/units/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UnitAdminRoutes(api fiber.Router, db *gorm.DB) {
	unitCtrl := unitController.NewUnitController(db)
	unitNewsCtrl := unitController.NewUnitNewsController(db)

	unitRoutes := api.Group("/units",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola unit"),
			constants.TeacherAndAbove,
		),
	)
	unitRoutes.Post("/", unitCtrl.CreateUnit)
	unitRoutes.Put("/:id", unitCtrl.UpdateUnit)
	unitRoutes.Delete("/:id", unitCtrl.DeleteUnit)

	unitNewsRoutes := api.Group("/units-news",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola berita unit"),
			constants.TeacherAndAbove,
		),
	)
	unitNewsRoutes.Post("/", unitNewsCtrl.Create)
	unitNewsRoutes.Put("/:id", unitNewsCtrl.Update)
	unitNewsRoutes.Delete("/:id", unitNewsCtrl.Delete)
}
