package route

import (
	difficultyController "quizku/internals/features/lessons/difficulty/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DifficultyUserRoutes(api fiber.Router, db *gorm.DB) {
	difficultyCtrl := difficultyController.NewDifficultyController(db)
	difficultyNewsCtrl := difficultyController.NewDifficultyNewsController(db)

	// ✅ Semua user login boleh akses
	difficultyRoutes := api.Group("/difficulties")
	difficultyRoutes.Get("/", difficultyCtrl.GetDifficulties)
	difficultyRoutes.Get("/:id", difficultyCtrl.GetDifficulty)

	difficultyNewsRoutes := api.Group("/difficulties-news")
	difficultyNewsRoutes.Get("/", difficultyNewsCtrl.GetAllNews)
	difficultyNewsRoutes.Get("/:difficulty_id", difficultyNewsCtrl.GetNewsByDifficultyId)
	difficultyNewsRoutes.Get("/detail/:id", difficultyNewsCtrl.GetNewsByID)
}
