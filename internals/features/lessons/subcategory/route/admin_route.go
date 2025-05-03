package route


import (
	subcategoryController "quizku/internals/features/lessons/subcategory/controller"
	"quizku/internals/constants"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubcategoryAdminRoutes(api fiber.Router, db *gorm.DB) {
	subcategoryCtrl := subcategoryController.NewSubcategoryController(db)
	subcategoryNewsCtrl := subcategoryController.NewSubcategoryNewsController(db)

	// 🔒 Subcategory
	subcategoryRoutes := api.Group("/subcategories",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola subkategori"),
			constants.TeacherAndAbove,
		),
	)
	subcategoryRoutes.Post("/", subcategoryCtrl.CreateSubcategory)
	subcategoryRoutes.Put("/:id", subcategoryCtrl.UpdateSubcategory)
	subcategoryRoutes.Delete("/:id", subcategoryCtrl.DeleteSubcategory)

	// 🔒 Subcategory News
	subcategoryNewsRoutes := api.Group("/subcategories-news",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola berita subkategori"),
			constants.TeacherAndAbove,
		),
	)
	subcategoryNewsRoutes.Post("/", subcategoryNewsCtrl.Create)
	subcategoryNewsRoutes.Put("/:id", subcategoryNewsCtrl.Update)
	subcategoryNewsRoutes.Delete("/:id", subcategoryNewsCtrl.Delete)
}
