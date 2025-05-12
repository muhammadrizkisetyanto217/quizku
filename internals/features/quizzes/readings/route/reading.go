package route

import (
	"quizku/internals/constants"
	readingController "quizku/internals/features/quizzes/readings/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReadingsRoute(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 📖 Reading Routes
	readingCtrl := readingController.NewReadingController(db)
	readingRoutes := api.Group("/readings")

	// ✅ GET routes boleh untuk semua user login
	readingRoutes.Get("/", readingCtrl.GetReadings)
	readingRoutes.Get("/:id", readingCtrl.GetReading)
	readingRoutes.Get("/unit/:unitId", readingCtrl.GetReadingsByUnit)

	// 🔒 POST, PUT, DELETE: hanya untuk pengelola (teacher, admin, owner)
	readingRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan reading"),
		constants.TeacherAndAbove,
	), readingCtrl.CreateReading)

	readingRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit reading"),
		constants.TeacherAndAbove,
	), readingCtrl.UpdateReading)

	readingRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus reading"),
		constants.TeacherAndAbove,
	), readingCtrl.DeleteReading)

	// 🧠 Tooltips integration (GET semua, boleh user biasa)
	readingRoutes.Get("/:id/with-tooltips", readingCtrl.GetReadingWithTooltips)

	// 📘 User Reading Routes (semua user login boleh akses)
	userReadingCtrl := readingController.NewUserReadingController(db)
	userReadingRoutes := api.Group("/user-readings")
	userReadingRoutes.Post("/", userReadingCtrl.CreateUserReading)
	userReadingRoutes.Get("/:id", userReadingCtrl.GetAllUserReading)
	userReadingRoutes.Get("/user/:user_id", userReadingCtrl.GetByUserID)
}
