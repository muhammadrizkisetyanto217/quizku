package route

import (
	themes_or_levelsController "quizku/internals/features/lessons/themes_or_levels/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ThemesOrLevelsUserRoutes(api fiber.Router, db *gorm.DB) {
	themeOrLevelCtrl := themes_or_levelsController.NewThemeOrLevelController(db)
	themesNewsCtrl := themes_or_levelsController.NewThemesOrLevelsNewsController(db)
	userThemesCtrl := themes_or_levelsController.NewUserThemesController(db)

	themeOrLevelRoutes := api.Group("/themes-or-levels")
	themeOrLevelRoutes.Get("/", themeOrLevelCtrl.GetThemeOrLevels)
	themeOrLevelRoutes.Get("/:id", themeOrLevelCtrl.GetThemeOrLevelById)
	themeOrLevelRoutes.Get("/subcategories/:subcategory_id", themeOrLevelCtrl.GetThemesOrLevelsBySubcategory)

	themesNewsRoutes := api.Group("/themes-or-levels-news")
	themesNewsRoutes.Get("/", themesNewsCtrl.GetAll)
	themesNewsRoutes.Get("/:themes_or_levels_id", themesNewsCtrl.GetByThemesOrLevelsID)
	themesNewsRoutes.Get("/detail/:id", themesNewsCtrl.GetByID)

	userThemesRoutes := api.Group("/user-themes-or-levels")
	userThemesRoutes.Get("/:user_id", userThemesCtrl.GetByUserID) // idealnya cek token vs ID
}