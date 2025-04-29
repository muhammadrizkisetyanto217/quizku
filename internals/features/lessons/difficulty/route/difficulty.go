package route

import (
	difficultyController "quizku/internals/features/lessons/difficulty/controller"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Group API: dilindungi Auth + Global RateLimiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 🎯 Difficulty Controller
	difficultyCtrl := difficultyController.NewDifficultyController(db)

	// 🔹 Group /difficulties
	difficultyRoutes := api.Group("/difficulties")
	difficultyRoutes.Get("/", difficultyCtrl.GetDifficulties)  // ✅ semua user login boleh
	difficultyRoutes.Get("/:id", difficultyCtrl.GetDifficulty) // ✅ semua user login boleh

	// 🔹 Untuk Create/Update/Delete baru butuh role tertentu
	protectedDifficultyRoutes := difficultyRoutes.Group("/",
		authMiddleware.OnlyRoles("❌ Hanya admin, teacher, atau owner yang bisa mengelola difficulties.", "admin", "teacher", "owner"),
	)
	protectedDifficultyRoutes.Post("/", difficultyCtrl.CreateDifficulty)
	protectedDifficultyRoutes.Put("/:id", difficultyCtrl.UpdateDifficulty)
	protectedDifficultyRoutes.Delete("/:id", difficultyCtrl.DeleteDifficulty)

	// 🎯 Difficulty News Controller
	difficultyNewsCtrl := difficultyController.NewDifficultyNewsController(db)

	// 🔹 Group /difficulties-news
	difficultyNewsRoutes := api.Group("/difficulties-news")
	difficultyNewsRoutes.Get("/:difficulty_id", difficultyNewsCtrl.GetNewsByDifficulty) // ✅ semua user login boleh
	difficultyNewsRoutes.Get("/detail/:id", difficultyNewsCtrl.GetNewsByID)             // ✅ semua user login boleh

	// 🔹 Untuk Create/Update/Delete News baru perlu role
	protectedDifficultyNewsRoutes := difficultyNewsRoutes.Group("/",
		authMiddleware.OnlyRoles("❌ Hanya admin, teacher, atau owner yang bisa mengelola difficulty news.", "admin", "teacher", "owner"),
	)
	protectedDifficultyNewsRoutes.Post("/", difficultyNewsCtrl.CreateNews)
	protectedDifficultyNewsRoutes.Put("/:id", difficultyNewsCtrl.UpdateNews)
	protectedDifficultyNewsRoutes.Delete("/:id", difficultyNewsCtrl.DeleteNews)
}
