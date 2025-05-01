package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	quizzesRoute "quizku/internals/features/quizzes/quizzes/route"
	evaluationRoute "quizku/internals/features/quizzes/evaluations/route"
	examsRoute "quizku/internals/features/quizzes/exams/route"
	readingsRoute "quizku/internals/features/quizzes/readings/route"
)

func QuizzesRoute(app *fiber.App, db *gorm.DB) {

	quizzesRoute.QuizzesRoutes(app, db)
	evaluationRoute.EvaluationsRoute(app, db)
	examsRoute.ExamsRoute(app, db)
	readingsRoute.ReadingsRoute(app, db)


}