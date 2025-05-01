package route

import (
	"quizku/internals/constants"
	subcategoryController "quizku/internals/features/lessons/subcategory/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🎯 Subcategory Routes
	subcategoryCtrl := subcategoryController.NewSubcategoryController(db)
	subcategoryRoutes := api.Group("/subcategories")

	// ✅ GET routes boleh semua login user
	subcategoryRoutes.Get("/", subcategoryCtrl.GetSubcategories)
	subcategoryRoutes.Get("/:id", subcategoryCtrl.GetSubcategory)
	subcategoryRoutes.Get("/category/:category_id", subcategoryCtrl.GetSubcategoriesByCategory)
	subcategoryRoutes.Get("/with-category-themes/:difficulty_id", subcategoryCtrl.GetCategoryWithSubcategoryAndThemes)

	// 🔒 CRUD dibatasi hanya untuk teacher/admin/owner
	subcategoryRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan subkategori"),
		constants.TeacherAndAbove,
	), subcategoryCtrl.CreateSubcategory)

	subcategoryRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit subkategori"),
		constants.TeacherAndAbove,
	), subcategoryCtrl.UpdateSubcategory)

	subcategoryRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus subkategori"),
		constants.TeacherAndAbove,
	), subcategoryCtrl.DeleteSubcategory)

	// 📰 Subcategory News Routes
	subcategoryNewsCtrl := subcategoryController.NewSubcategoryNewsController(db)
	subcategoryNewsRoutes := api.Group("/subcategories-news")

	subcategoryNewsRoutes.Get("/:subcategory_id", subcategoryNewsCtrl.GetAll)
	subcategoryNewsRoutes.Get("/:id", subcategoryNewsCtrl.GetByID)

	// 🔒 CRUD news juga hanya untuk teacher/admin/owner
	subcategoryNewsRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan berita subkategori"),
		constants.TeacherAndAbove,
	), subcategoryNewsCtrl.Create)

	subcategoryNewsRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit berita subkategori"),
		constants.TeacherAndAbove,
	), subcategoryNewsCtrl.Update)

	subcategoryNewsRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus berita subkategori"),
		constants.TeacherAndAbove,
	), subcategoryNewsCtrl.Delete)

	// ✅ User Subcategory Routes → semua login user
	userSubcategoryCtrl := subcategoryController.NewUserSubcategoryController(db)
	userSubcategoryRoutes := api.Group("/user-subcategory")
	userSubcategoryRoutes.Post("/", userSubcategoryCtrl.Create)
	userSubcategoryRoutes.Get("/:id", userSubcategoryCtrl.GetByUserId)
	userSubcategoryRoutes.Get("/category/:user_id/difficulty/:difficulty_id", userSubcategoryCtrl.GetWithProgressByParam)
}
