package category

import (
	"quizku/internals/constants"
	categoryController "quizku/internals/features/lessons/categories/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	// 🎯 Category Routes
	categoryCtrl := categoryController.NewCategoryController(db)
	categoryRoutes := api.Group("/categories")

	// ✅ GET kategori bisa diakses semua user login
	categoryRoutes.Get("/", categoryCtrl.GetCategories)
	categoryRoutes.Get("/:id", categoryCtrl.GetCategory)
	categoryRoutes.Get("/difficulty/:difficulty_id", categoryCtrl.GetCategoriesByDifficulty)

	// 🔒 Hanya pengelola yang bisa CRUD kategori
	categoryRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan kategori"),
		constants.TeacherAndAbove,
	), categoryCtrl.CreateCategory)

	categoryRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit kategori"),
		constants.TeacherAndAbove,
	), categoryCtrl.UpdateCategory)

	categoryRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus kategori"),
		constants.TeacherAndAbove,
	), categoryCtrl.DeleteCategory)

	// 📰 Category News Routes
	categoryNewsCtrl := categoryController.NewCategoryNewsController(db)
	categoryNewsRoutes := api.Group("/categories-news")

	// ✅ GET berita kategori terbuka untuk user login
	categoryNewsRoutes.Get("/:category_id", categoryNewsCtrl.GetAll)
	categoryNewsRoutes.Get("/:id", categoryNewsCtrl.GetByID)

	// 🔒 CRUD hanya untuk pengelola
	categoryNewsRoutes.Post("/", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menambahkan berita kategori"),
		constants.TeacherAndAbove,
	), categoryNewsCtrl.Create)

	categoryNewsRoutes.Put("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("mengedit berita kategori"),
		constants.TeacherAndAbove,
	), categoryNewsCtrl.Update)

	categoryNewsRoutes.Delete("/:id", authMiddleware.OnlyRolesSlice(
		constants.RoleErrorTeacher("menghapus berita kategori"),
		constants.TeacherAndAbove,
	), categoryNewsCtrl.Delete)
}
