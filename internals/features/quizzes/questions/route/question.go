package route

import (
	"quizku/internals/constants"
	"quizku/internals/features/quizzes/questions/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuestionRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🎯 Quiz Question Routes
	questionController := controller.NewQuestionController(db)
	questionRoutes := api.Group("/question")

	// ✅ GET routes boleh semua user login
	questionRoutes.Get("/", questionController.GetQuestions)
	questionRoutes.Get("/:id", questionController.GetQuestion)
	questionRoutes.Get("/:quizId/questionsQuiz", questionController.GetQuestionsByQuizID)
	questionRoutes.Get("/:evaluationId/questionsEvaluation", questionController.GetQuestionsByEvaluationID)
	questionRoutes.Get("/:examId/questionsExam", questionController.GetQuestionsByExamID)
	questionRoutes.Get("/:id/questionTooltips", questionController.GetQuestionWithTooltips)
	questionRoutes.Get("/:id/questionTooltips/:tooltipId", questionController.GetOnlyQuestionTooltips)
	questionRoutes.Get("/:id/questionTooltipsMarked", questionController.GetQuestionWithTooltipsMarked)

	// 🔒 POST, PUT, DELETE hanya untuk teacher/admin/owner
	questionRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("membuat soal"),
		constants.TeacherAndAbove,
	), questionController.CreateQuestion)

	questionRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit soal"),
		constants.TeacherAndAbove,
	), questionController.UpdateQuestion)

	questionRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus soal"),
		constants.TeacherAndAbove,
	), questionController.DeleteQuestion)

	// ✅ Quiz Saved Routes (semua user login)
	questionSavedController := controller.NewQuestionSavedController(db)
	questionSavedRoutes := api.Group("/question-saved")
	questionSavedRoutes.Post("/", questionSavedController.Create)
	questionSavedRoutes.Get("/user/:user_id", questionSavedController.GetByUserID)
	questionSavedRoutes.Get("/question_saved_with_question/:user_id", questionSavedController.GetByUserIDWithQuestions)
	questionSavedRoutes.Delete("/user/:id", questionSavedController.Delete)

	// ✅ Quiz Mistake Routes (semua user login)
	questionMistakeController := controller.NewQuestionMistakeController(db)
	questionMistakeRoutes := api.Group("/question-mistakes")
	questionMistakeRoutes.Post("/", questionMistakeController.Create)
	questionMistakeRoutes.Get("/user/:user_id", questionMistakeController.GetByUserID)
	questionMistakeRoutes.Delete("/:id", questionMistakeController.Delete)

	// ✅ User Question Routes (semua user login)
	userQuestionController := controller.NewUserQuestionController(db)
	userQuestionRoutes := api.Group("/user-questions")
	userQuestionRoutes.Post("/", userQuestionController.Create)
	userQuestionRoutes.Get("/user/:user_id", userQuestionController.GetByUserID)
	userQuestionRoutes.Get("/user/:user_id/question/:question_id", userQuestionController.GetByUserIDAndQuestionID)
}
