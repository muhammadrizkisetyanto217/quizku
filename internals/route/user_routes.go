package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	userRoute "quizku/internals/features/users/user/routes"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 🔓 Prefix user biasa: /api/u/...
	userGroup := api.Group("/u") // 👤 user login biasa
	userRoute.UserAllRoutes(userGroup, db)

	// 🔐 Prefix admin: /api/a/...
	adminGroup := api.Group("/a") // 🔐 hanya teacher/admin/owner
	userRoute.UserAdminRoutes(adminGroup, db)

}
