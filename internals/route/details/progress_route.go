package details

import (
	levelRankRoute "quizku/internals/features/progress/level_rank/route"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProgressRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	adminGroup := api.Group("/a")
	levelRankRoute.LevelRequirementAdminRoute(adminGroup, db)

	userGroup := api.Group("/u")
	levelRankRoute.LevelRequirementUserRoute(userGroup, db)
}
