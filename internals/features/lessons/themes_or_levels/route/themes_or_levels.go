package route

import (
	"quizku/internals/constants"
	themes_or_levelsController "quizku/internals/features/lessons/themes_or_levels/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ThemeOrLevelRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🎯 Themes or Levels Routes
	themeOrLevelCtrl := themes_or_levelsController.NewThemeOrLevelController(db)
	themeOrLevelRoutes := api.Group("/themes-or-levels")

	// ✅ GET: bebas untuk semua user login
	themeOrLevelRoutes.Get("/", themeOrLevelCtrl.GetThemeOrLevels)
	themeOrLevelRoutes.Get("/:id", themeOrLevelCtrl.GetThemeOrLevelById)
	themeOrLevelRoutes.Get("/subcategories/:subcategory_id", themeOrLevelCtrl.GetThemesOrLevelsBySubcategory)

	// 🔒 POST, PUT, DELETE: hanya untuk teacher/admin/owner
	themeOrLevelRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan tema atau level"),
		constants.TeacherAndAbove,
	), themeOrLevelCtrl.CreateThemeOrLevel)

	themeOrLevelRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit tema atau level"),
		constants.TeacherAndAbove,
	), themeOrLevelCtrl.UpdateThemeOrLevel)

	themeOrLevelRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus tema atau level"),
		constants.TeacherAndAbove,
	), themeOrLevelCtrl.DeleteThemeOrLevel)

	// 📰 Themes or Levels News Routes
	themesNewsCtrl := themes_or_levelsController.NewThemesOrLevelsNewsController(db)
	themesNewsRoutes := api.Group("/themes-or-levels-news")

	// ✅ GET: semua user login
	themesNewsRoutes.Get("/", themesNewsCtrl.GetAll)
	themesNewsRoutes.Get("/:id", themesNewsCtrl.GetByID)

	// 🔒 POST, PUT, DELETE: hanya untuk pengelola
	themesNewsRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan berita tema atau level"),
		constants.TeacherAndAbove,
	), themesNewsCtrl.Create)

	themesNewsRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit berita tema atau level"),
		constants.TeacherAndAbove,
	), themesNewsCtrl.Update)

	themesNewsRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus berita tema atau level"),
		constants.TeacherAndAbove,
	), themesNewsCtrl.Delete)

	// ✅ User Themes or Levels Route
	userThemesCtrl := themes_or_levelsController.NewUserThemesController(db)
	userThemesRoutes := api.Group("/user-themes-or-levels")
	userThemesRoutes.Get("/:user_id", userThemesCtrl.GetByUserID) // idealnya validasi user_id vs user login
}
