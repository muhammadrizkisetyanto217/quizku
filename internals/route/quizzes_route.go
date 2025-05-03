package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	evaluationRoute "quizku/internals/features/quizzes/evaluations/route"
	examsRoute "quizku/internals/features/quizzes/exams/route"
	quizzesRoute "quizku/internals/features/quizzes/quizzes/route"
	readingsRoute "quizku/internals/features/quizzes/readings/route"

	questionsRoute "quizku/internals/features/quizzes/questions/route"


)

func QuizzesRoute(app *fiber.App, db *gorm.DB) {
	// 🔐 Bungkus dengan Auth dan RateLimiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	quizzesRoute.QuizzesAdminRoutes(api, db)
	quizzesRoute.QuizzesUserRoutes(api, db)

	evaluationRoute.EvaluationAdminRoutes(api, db)
	evaluationRoute.EvaluationUserRoutes(api, db)

	examsRoute.ExamAdminRoutes(api, db)
	examsRoute.ExamUserRoutes(api, db)

	readingsRoute.ReadingAdminRoutes(api, db)
	readingsRoute.ReadingUserRoutes(api, db)

	questionsRoute.QuestionAdminRoutes(api, db)
	questionsRoute.QuestionUserRoutes(api, db)



}
