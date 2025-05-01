package route

import (
	"quizku/internals/constants"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	levelController "quizku/internals/features/progress/level_rank/controller"
	authMiddleware "quizku/internals/middlewares/auth"
)

func LevelRequirementRoute(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	levelCtrl := levelController.NewLevelRequirementController(db)
	rankCtrl := levelController.NewRankRequirementController(db)

	// 🎯 Level Routes
	levelRoutes := api.Group("/level-requirements")
	levelRoutes.Get("/", levelCtrl.GetAll)
	levelRoutes.Get("/:id", levelCtrl.GetByID)

	levelRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan level"),
		constants.TeacherAndAbove,
	), levelCtrl.Create)

	levelRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit level"),
		constants.TeacherAndAbove,
	), levelCtrl.Update)

	levelRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus level"),
		constants.TeacherAndAbove,
	), levelCtrl.Delete)

	// 🏆 Rank Routes
	rankRoutes := api.Group("/rank-requirements")
	rankRoutes.Get("/", rankCtrl.GetAll)
	rankRoutes.Get("/:id", rankCtrl.GetByID)

	rankRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan rank"),
		constants.TeacherAndAbove,
	), rankCtrl.Create)

	rankRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit rank"),
		constants.TeacherAndAbove,
	), rankCtrl.Update)

	rankRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus rank"),
		constants.TeacherAndAbove,
	), rankCtrl.Delete)
}
