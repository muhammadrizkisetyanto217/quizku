package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	difficultiesRoute "quizku/internals/features/lessons/difficulty/route"
	categoriesRoute "quizku/internals/features/lessons/categories/route"
	subcategoriesRoute "quizku/internals/features/lessons/subcategory/route"
	themesOrLevelsRoute "quizku/internals/features/lessons/themes_or_levels/route"
	unitsRoute "quizku/internals/features/lessons/units/route"
)

func LessonRoutes(app *fiber.App, db *gorm.DB) {
	
	difficultiesRoute.DifficultyRoutes(app, db)
	categoriesRoute.CategoryRoutes(app, db)
	subcategoriesRoute.SubategoryRoutes(app, db)
	themesOrLevelsRoute.ThemesOrLevelsRoutes(app, db)
	unitsRoute.UnitRoutes(app, db)

}
