package route


import (
	subcategoryController "quizku/internals/features/lessons/subcategory/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubcategoryUserRoutes(api fiber.Router, db *gorm.DB) {
	subcategoryCtrl := subcategoryController.NewSubcategoryController(db)
	subcategoryNewsCtrl := subcategoryController.NewSubcategoryNewsController(db)
	userSubcategoryCtrl := subcategoryController.NewUserSubcategoryController(db)

	subcategoryRoutes := api.Group("/subcategories")
	subcategoryRoutes.Get("/", subcategoryCtrl.GetSubcategories)
	subcategoryRoutes.Get("/:id", subcategoryCtrl.GetSubcategory)
	subcategoryRoutes.Get("/category/:category_id", subcategoryCtrl.GetSubcategoriesByCategory)
	subcategoryRoutes.Get("/with-category-themes/:difficulty_id", subcategoryCtrl.GetCategoryWithSubcategoryAndThemes)

	subcategoryNewsRoutes := api.Group("/subcategories-news")
	subcategoryNewsRoutes.Get("/:subcategory_id", subcategoryNewsCtrl.GetAll)
	subcategoryNewsRoutes.Get("/:id", subcategoryNewsCtrl.GetByID)

	userSubcategoryRoutes := api.Group("/user-subcategory")
	userSubcategoryRoutes.Post("/", userSubcategoryCtrl.Create)
	userSubcategoryRoutes.Get("/:id", userSubcategoryCtrl.GetByUserId)
	userSubcategoryRoutes.Get("/category/difficulty/:difficulty_id", userSubcategoryCtrl.GetWithProgressByParam)
}
