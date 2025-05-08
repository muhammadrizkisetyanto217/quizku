package route

import (
	"quizku/internals/constants"
	testExamController "quizku/internals/features/users/test_exam/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestExamAdminRoutes(app fiber.Router, db *gorm.DB) {
	testExamCtrl := testExamController.NewTestExamController(db)
	userExamCtrl := testExamController.NewUserTestExamController(db)

	// 🔐 /test-exams – hanya untuk adminOwner
	adminTestExams := app.Group("/test-exams",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("Kelola Ujian"),
			constants.OwnerAndAbove,
		),
	)

	adminTestExams.Post("/", testExamCtrl.Create)
	adminTestExams.Put("/:id", testExamCtrl.Update)
	adminTestExams.Delete("/:id", testExamCtrl.Delete)

	// 🔐 /user-test-exam – hanya untuk adminOwner
	userTestExams := app.Group("/user-test-exams",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("Kelola User Test Exam"),
			constants.OwnerAndAbove,
		),
	)

	userTestExams.Get("/", userExamCtrl.GetAll)
	userTestExams.Get("/exam/:test_exam_id", userExamCtrl.GetByTestExamID)
}
