package route

import (
	"quizku/internals/constants"
	examController "quizku/internals/features/quizzes/exams/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExamRoute(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 📝 Exam Routes
	examCtrl := examController.NewExamController(db)
	examRoutes := api.Group("/exams")

	// ✅ GET exam bebas diakses semua user login
	examRoutes.Get("/", examCtrl.GetExams)
	examRoutes.Get("/:id", examCtrl.GetExam)
	examRoutes.Get("/unit/:unitId", examCtrl.GetExamsByUnitID)

	// 🔒 POST, PUT, DELETE exam hanya untuk teacher/admin/owner
	examRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("membuat ujian"),
		constants.TeacherAndAbove,
	), examCtrl.CreateExam)

	examRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit ujian"),
		constants.TeacherAndAbove,
	), examCtrl.UpdateExam)

	examRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus ujian"),
		constants.TeacherAndAbove,
	), examCtrl.DeleteExam)

	// ✅ User Exam Routes (semua user login)
	userExamCtrl := examController.NewUserExamController(db)
	userExamRoutes := api.Group("/user-exams")

	userExamRoutes.Post("/", userExamCtrl.Create)
	userExamRoutes.Get("/", userExamCtrl.GetAll)
	userExamRoutes.Get("/user/:user_id", userExamCtrl.GetByUserID)
	userExamRoutes.Get("/:id", userExamCtrl.GetByID)
	userExamRoutes.Delete("/:id", userExamCtrl.Delete)
}
