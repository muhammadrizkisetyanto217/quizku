package route

import (
	userController "quizku/internals/features/users/user/controller"
	rateLimiter "quizku/internals/middlewares" //
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mengatur routing untuk user & user profile
func UserRoutes(app *fiber.App, db *gorm.DB) {

	// ✅ Group API umum, dilindungi Auth + Global RateLimiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db), // Auth dulu
		rateLimiter.GlobalRateLimiter(),   // Baru RateLimiter global
	)

	// 🔹 Users Controller
	userCtrl := userController.NewUserController(db)

	// 🔹 Group /users: hanya untuk teacher dan owner
	userRoutes := api.Group("/users",
		authMiddleware.OnlyRoles("❌ Hanya teacher atau owner yang bisa mengakses user management.", "teacher", "owner"),
	)
	userRoutes.Get("/", userCtrl.GetUsers)
	userRoutes.Get("/profile", userCtrl.GetProfile)
	userRoutes.Put("/profile", userCtrl.UpdateProfile)
	userRoutes.Delete("/:id", userCtrl.DeleteUser)

	// 🔹 Users Profile: semua user yang sudah login boleh akses
	userProfileCtrl := userController.NewUsersProfileController(db)
	usersProfileRoutes := api.Group("/users-profiles")
	usersProfileRoutes.Get("/", userProfileCtrl.GetProfiles)
	usersProfileRoutes.Get("/:id", userProfileCtrl.GetProfile)
	usersProfileRoutes.Post("/", userProfileCtrl.CreateProfile)
	usersProfileRoutes.Put("/:id", userProfileCtrl.UpdateProfile)
	usersProfileRoutes.Delete("/:id", userProfileCtrl.DeleteProfile)
}
