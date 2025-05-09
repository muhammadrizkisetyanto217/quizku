package route

import (
	"quizku/internals/features/quizzes/questions/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuestionUserRoutes(api fiber.Router, db *gorm.DB) {
	questionController := controller.NewQuestionController(db)
	questionSavedController := controller.NewQuestionSavedController(db)
	questionMistakeController := controller.NewQuestionMistakeController(db)
	userQuestionController := controller.NewUserQuestionController(db)

	questionRoutes := api.Group("/question")
	questionRoutes.Get("/", questionController.GetQuestions)
	questionRoutes.Get("/:id", questionController.GetQuestion)

	//* Main
	questionRoutes.Get("/:quizId/questionsQuiz", questionController.GetQuestionsByQuizID)
	questionRoutes.Get("/:evaluationId/questionsEvaluation", questionController.GetQuestionsByEvaluationID)
	questionRoutes.Get("/:examId/questionsExam", questionController.GetQuestionsByExamID)
	questionRoutes.Get("/:testId/questionsTest", questionController.GetQuestionsByTestID)
	// 

	questionRoutes.Get("/:id/questionTooltips", questionController.GetQuestionWithTooltips)
	questionRoutes.Get("/:id/questionTooltips/:tooltipId", questionController.GetOnlyQuestionTooltips)
	questionRoutes.Get("/:id/questionTooltipsMarked", questionController.GetQuestionWithTooltipsMarked)

	// Question Saved
	questionSavedRoutes := api.Group("/question-saved")
	questionSavedRoutes.Post("/", questionSavedController.Create)
	questionSavedRoutes.Get("/user/:user_id", questionSavedController.GetByUserID)
	questionSavedRoutes.Get("/question_saved_with_question/:user_id", questionSavedController.GetByUserIDWithQuestions)
	questionSavedRoutes.Delete("/user/:id", questionSavedController.Delete)

	// Question Mistakes
	questionMistakeRoutes := api.Group("/question-mistakes")
	questionMistakeRoutes.Post("/", questionMistakeController.Create)
	questionMistakeRoutes.Get("/user/:user_id", questionMistakeController.GetByUserID)
	questionMistakeRoutes.Delete("/:id", questionMistakeController.Delete)

	// User Question Progress
	userQuestionRoutes := api.Group("/user-questions")
	userQuestionRoutes.Post("/", userQuestionController.Create)
	userQuestionRoutes.Get("/user/:user_id", userQuestionController.GetByUserID)
	userQuestionRoutes.Get("/user/:user_id/question/:question_id", userQuestionController.GetByUserIDAndQuestionID)
}
