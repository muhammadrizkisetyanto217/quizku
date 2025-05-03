package details

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	evaluationRoute "quizku/internals/features/quizzes/evaluations/route"
	examsRoute "quizku/internals/features/quizzes/exams/route"
	questionsRoute "quizku/internals/features/quizzes/questions/route"
	quizzesRoute "quizku/internals/features/quizzes/quizzes/route"
	readingsRoute "quizku/internals/features/quizzes/readings/route"
)

func QuizzesRoute(app *fiber.App, db *gorm.DB) {
	// Bungkus dengan Auth dan RateLimiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 👤 Prefix user: /api/u
	userGroup := api.Group("/u")
	quizzesRoute.QuizzesUserRoutes(userGroup, db)
	evaluationRoute.EvaluationUserRoutes(userGroup, db)
	examsRoute.ExamUserRoutes(userGroup, db)
	readingsRoute.ReadingUserRoutes(userGroup, db)
	questionsRoute.QuestionUserRoutes(userGroup, db)

	// 🔐 Prefix admin: /api/a
	adminGroup := api.Group("/a")
	quizzesRoute.QuizzesAdminRoutes(adminGroup, db)
	evaluationRoute.EvaluationAdminRoutes(adminGroup, db)
	examsRoute.ExamAdminRoutes(adminGroup, db)
	readingsRoute.ReadingAdminRoutes(adminGroup, db)
	questionsRoute.QuestionAdminRoutes(adminGroup, db)
}
