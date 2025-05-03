package route

import (
	categoryController "quizku/internals/features/lessons/categories/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryUserRoutes(api fiber.Router, db *gorm.DB) {
	categoryCtrl := categoryController.NewCategoryController(db)
	categoryNewsCtrl := categoryController.NewCategoryNewsController(db)

	categoryRoutes := api.Group("/categories")
	categoryRoutes.Get("/", categoryCtrl.GetCategories)
	categoryRoutes.Get("/:id", categoryCtrl.GetCategory)
	categoryRoutes.Get("/difficulty/:difficulty_id", categoryCtrl.GetCategoriesByDifficulty)

	categoryNewsRoutes := api.Group("/categories-news")
	categoryNewsRoutes.Get("/:category_id", categoryNewsCtrl.GetAll)
	categoryNewsRoutes.Get("/:id", categoryNewsCtrl.GetByID)
}
