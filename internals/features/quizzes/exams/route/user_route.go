package route

import (
	examController "quizku/internals/features/quizzes/exams/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExamUserRoutes(api fiber.Router, db *gorm.DB) {
	examCtrl := examController.NewExamController(db)
	userExamCtrl := examController.NewUserExamController(db)

	examRoutes := api.Group("/exams")
	examRoutes.Get("/", examCtrl.GetExams)
	examRoutes.Get("/:id", examCtrl.GetExam)
	examRoutes.Get("/unit/:unitId", examCtrl.GetExamsByUnitID)

	userExamRoutes := api.Group("/user-exams")
	userExamRoutes.Post("/", userExamCtrl.Create)
	userExamRoutes.Get("/", userExamCtrl.GetAll)
	userExamRoutes.Get("/user/:user_id", userExamCtrl.GetByUserID)
	userExamRoutes.Get("/:id", userExamCtrl.GetByID)
	userExamRoutes.Delete("/:id", userExamCtrl.Delete)
}
