package route

import (
	themes_or_levelsController "quizku/internals/features/lessons/themes_or_levels/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ThemesOrLevelsAdminRoutes(api fiber.Router, db *gorm.DB) {
	themeOrLevelCtrl := themes_or_levelsController.NewThemeOrLevelController(db)
	themesNewsCtrl := themes_or_levelsController.NewThemesOrLevelsNewsController(db)

	themeOrLevelRoutes := api.Group("/themes-or-levels",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola tema atau level"),
			constants.OwnerAndAbove,
		),
	)
	themeOrLevelRoutes.Post("/", themeOrLevelCtrl.CreateThemeOrLevel)
	themeOrLevelRoutes.Put("/:id", themeOrLevelCtrl.UpdateThemeOrLevel)
	themeOrLevelRoutes.Delete("/:id", themeOrLevelCtrl.DeleteThemeOrLevel)

	themesNewsRoutes := api.Group("/themes-or-levels-news",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola berita tema atau level"),
			constants.OwnerAndAbove,
		),
	)
	themesNewsRoutes.Post("/", themesNewsCtrl.Create)
	themesNewsRoutes.Put("/:id", themesNewsCtrl.Update)
	themesNewsRoutes.Delete("/:id", themesNewsCtrl.Delete)
}