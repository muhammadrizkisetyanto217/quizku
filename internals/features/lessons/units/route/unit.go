package route

import (
	"quizku/internals/constants"
	unitController "quizku/internals/features/lessons/units/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UnitRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🎯 Unit Routes
	unitCtrl := unitController.NewUnitController(db)
	unitRoutes := api.Group("/units")

	// ✅ GET (semua user login)
	unitRoutes.Get("/", unitCtrl.GetUnits)
	unitRoutes.Get("/:id", unitCtrl.GetUnit)
	unitRoutes.Get("/themes-or-levels/:themesOrLevelId", unitCtrl.GetUnitByThemesOrLevels)

	// 🔒 POST/PUT/DELETE → hanya pengelola
	unitRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan unit"),
		constants.TeacherAndAbove,
	), unitCtrl.CreateUnit)

	unitRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit unit"),
		constants.TeacherAndAbove,
	), unitCtrl.UpdateUnit)

	unitRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus unit"),
		constants.TeacherAndAbove,
	), unitCtrl.DeleteUnit)

	// 📰 Unit News Routes
	unitNewsCtrl := unitController.NewUnitNewsController(db)
	unitNewsRoutes := api.Group("/units-news")

	// ✅ GET (semua user login)
	unitNewsRoutes.Get("/", unitNewsCtrl.GetAll)
	unitNewsRoutes.Get("/:id", unitNewsCtrl.GetByID)

	// 🔒 POST/PUT/DELETE → hanya pengelola
	unitNewsRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan berita unit"),
		constants.TeacherAndAbove,
	), unitNewsCtrl.Create)

	unitNewsRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit berita unit"),
		constants.TeacherAndAbove,
	), unitNewsCtrl.Update)

	unitNewsRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus berita unit"),
		constants.TeacherAndAbove,
	), unitNewsCtrl.Delete)

	// 👤 User Unit Routes
	userUnitCtrl := unitController.NewUserUnitController(db)
	userUnitRoutes := api.Group("/user-units")
	userUnitRoutes.Get("/:user_id", userUnitCtrl.GetByUserID)
	userUnitRoutes.Get("/:user_id/themes-or-levels/:themes_or_levels_id", userUnitCtrl.GetUserUnitsByThemesOrLevelsAndUserID)
}
