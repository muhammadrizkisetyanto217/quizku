package route

import (
	"log"
	testExamController "quizku/internals/features/users/test_exam/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestExamUserRoutes(app fiber.Router, db *gorm.DB) {
	log.Println("[DEBUG] ❗ Masuk TestExamAllRoutes")

	testExamCtrl := testExamController.NewTestExamController(db)

	// ✅ Bisa diakses semua user login
	app.Get("/test-exams", testExamCtrl.GetAll)
	app.Get("/test-exams/:id", testExamCtrl.GetByID)

	userExamCtrl := testExamController.NewUserTestExamController(db)
	
	// ✅ Akses user login biasa
	app.Get("/user-test-exams/user/:user_id", userExamCtrl.GetByUserID)
	app.Get("/user-test-exams/:id", userExamCtrl.GetByID)
	app.Post("/user-test-exams", userExamCtrl.Create)
}
