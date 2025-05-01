package route

import (
	"quizku/internals/constants"
	quizzesController "quizku/internals/features/quizzes/quizzes/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuizzesRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🔥 Section Quizzes Routes
	sectionQuizzesController := quizzesController.NewSectionQuizController(db)
	sectionQuizzesRoutes := api.Group("/section-quizzes")

	// ✅ GET bisa diakses semua user login
	sectionQuizzesRoutes.Get("/", sectionQuizzesController.GetSectionQuizzes)
	sectionQuizzesRoutes.Get("/:id", sectionQuizzesController.GetSectionQuiz)
	sectionQuizzesRoutes.Get("/unit/:unitId", sectionQuizzesController.GetSectionQuizzesByUnit)

	// 🔒 Create, Update, Delete hanya untuk teacher/admin/owner
	sectionQuizzesRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("membuat section quiz"),
		constants.TeacherAndAbove,
	), sectionQuizzesController.CreateSectionQuiz)

	sectionQuizzesRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengubah section quiz"),
		constants.TeacherAndAbove,
	), sectionQuizzesController.UpdateSectionQuiz)

	sectionQuizzesRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus section quiz"),
		constants.TeacherAndAbove,
	), sectionQuizzesController.DeleteSectionQuiz)

	// 🧠 Quiz Routes
	quizController := quizzesController.NewQuizController(db)
	quizRoutes := api.Group("/quizzes")

	// ✅ GET quiz bebas untuk semua user login
	quizRoutes.Get("/", quizController.GetQuizzes)
	quizRoutes.Get("/:id", quizController.GetQuiz)
	quizRoutes.Get("/section/:sectionId", quizController.GetQuizzesBySection)

	// 🔒 POST, PUT, DELETE quiz hanya untuk teacher/admin/owner
	quizRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("membuat quiz"),
		constants.TeacherAndAbove,
	), quizController.CreateQuiz)

	quizRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengubah quiz"),
		constants.TeacherAndAbove,
	), quizController.UpdateQuiz)

	quizRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus quiz"),
		constants.TeacherAndAbove,
	), quizController.DeleteQuiz)

	// 🧑‍🎓 User Quiz Routes (semua user login)
	userQuizController := quizzesController.NewUserQuizController(db)
	userQuizRoutes := api.Group("/user-quizzes")
	userQuizRoutes.Post("/", userQuizController.CreateOrUpdateUserQuiz)
	userQuizRoutes.Get("/user/:user_id", userQuizController.GetUserQuizzesByUserID)

	// 🧑‍🎓 User Section Quizzes Routes (semua user login)
	userSectionQuizzesController := quizzesController.NewUserSectionQuizzesController(db)
	userSectionQuizzesRoutes := api.Group("/user-section-quizzes")
	userSectionQuizzesRoutes.Get("/user/:user_id", userSectionQuizzesController.GetUserSectionQuizzesByUserID)
}
