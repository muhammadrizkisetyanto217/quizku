package routes

import (
	categoriesRoute "quizku/internals/features/lessons/categories/route"
	difficultiesRoute "quizku/internals/features/lessons/difficulty/route"
	subcategoriesRoute "quizku/internals/features/lessons/subcategory/route"
	themesOrLevelsRoute "quizku/internals/features/lessons/themes_or_levels/route"
	unitsRoute "quizku/internals/features/lessons/units/route"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LessonRoutes(app *fiber.App, db *gorm.DB) {
	// ✅ Semua route lesson berada di bawah /api dan dilindungi auth + rate limiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// ✅ Route dibagi dua: user dan admin
	difficultiesRoute.DifficultyUserRoutes(api, db)
	difficultiesRoute.DifficultyAdminRoutes(api, db)
	
	categoriesRoute.CategoryAdminRoutes(api, db)
	categoriesRoute.CategoryUserRoutes(api, db)


	subcategoriesRoute.SubcategoryAdminRoutes(api, db)
	subcategoriesRoute.SubcategoryUserRoutes(api, db)

	themesOrLevelsRoute.ThemesOrLevelsAdminRoutes(api, db)
	themesOrLevelsRoute.ThemesOrLevelsUserRoutes(api, db)

	unitsRoute.UnitAdminRoutes(api, db)
	unitsRoute.UnitUserRoutes(api, db)

}
