package routes

import (
	"quizku/internals/constants" // ✅ Tambahkan ini
	userController "quizku/internals/features/users/user/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserAdminRoutes(app fiber.Router, db *gorm.DB) {
	userCtrl := userController.NewUserController(db)
	userProfileCtrl := userController.NewUsersProfileController(db)

	// 🔐 /users – hanya teacher, admin, owner
	users := app.Group("/users",
		authMiddleware.OnlyRolesSlice(constants.RoleErrorTeacher("User Management"), constants.TeacherAndAbove),
	)

	users.Get("/", userCtrl.GetUsers)
	users.Put("/user", userCtrl.UpdateUser)
	users.Post("/", userCtrl.CreateUser)
	users.Delete("/:id", userCtrl.DeleteUser)

	// 🔐 Tambahan: admin bisa lihat semua user profile
	app.Get("/users-profiles",
		authMiddleware.OnlyRolesSlice(constants.RoleErrorTeacher("Lihat Semua User Profile"), constants.TeacherAndAbove),
		userProfileCtrl.GetProfiles,
	)

}
