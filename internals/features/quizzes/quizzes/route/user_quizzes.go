package route

import (
	quizzesController "quizku/internals/features/quizzes/quizzes/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuizzesUserRoutes(api fiber.Router, db *gorm.DB) {
	sectionQuizzesController := quizzesController.NewSectionQuizController(db)
	quizController := quizzesController.NewQuizController(db)
	userQuizController := quizzesController.NewUserQuizController(db)
	userSectionQuizzesController := quizzesController.NewUserSectionQuizzesController(db)

	// Section Quizzes - GET only
	sectionRoutes := api.Group("/section-quizzes")
	sectionRoutes.Get("/", sectionQuizzesController.GetSectionQuizzes)
	sectionRoutes.Get("/:id", sectionQuizzesController.GetSectionQuiz)
	sectionRoutes.Get("/unit/:unitId", sectionQuizzesController.GetSectionQuizzesByUnit)

	// Quizzes - GET only
	quizRoutes := api.Group("/quizzes")
	quizRoutes.Get("/", quizController.GetQuizzes)
	quizRoutes.Get("/:id", quizController.GetQuiz)
	quizRoutes.Get("/section/:sectionId", quizController.GetQuizzesBySection)

	// User Quizzes
	userQuizRoutes := api.Group("/user-quizzes")
	userQuizRoutes.Post("/", userQuizController.CreateOrUpdateUserQuiz)
	userQuizRoutes.Get("/user/:user_id", userQuizController.GetUserQuizzesByUserID)

	// User Section Quizzes
	userSectionRoutes := api.Group("/user-section-quizzes")
	userSectionRoutes.Get("/user/:user_id", userSectionQuizzesController.GetUserSectionQuizzesByUserID)
}
