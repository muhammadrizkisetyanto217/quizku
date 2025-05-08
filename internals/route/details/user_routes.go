package details

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	surveyRoute "quizku/internals/features/users/survey/route"
	tokenRoute "quizku/internals/features/users/token/route"
	userRoute "quizku/internals/features/users/user/routes"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"
	testExamRoute "quizku/internals/features/users/test_exam/route"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	adminGroup := api.Group("/a") // 🔐 hanya teacher/admin/owner
	userRoute.UserAdminRoutes(adminGroup, db)
	surveyRoute.SurveyAdminRoutes(adminGroup, db)
	testExamRoute.TestExamAdminRoutes(adminGroup, db) // 🔐 Ujian Admin

	// 🔓 Prefix user biasa: /api/u/...
	userGroup := api.Group("/u") // 👤 user login biasa
	userRoute.UserAllRoutes(userGroup, db)
	surveyRoute.SurveyUserRoutes(userGroup, db)
	tokenRoute.RegisterTokenRoutes(userGroup, db) // 🔓 Token routes
	testExamRoute.TestExamUserRoutes(userGroup, db) // 👤 Ujian User

}
