package route

import (
	difficultyController "quizku/internals/features/lessons/difficulty/controller"
	authMiddleware "quizku/internals/middlewares/auth"
	"quizku/internals/constants"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DifficultyAdminRoutes(api fiber.Router, db *gorm.DB) {
	difficultyCtrl := difficultyController.NewDifficultyController(db)
	difficultyNewsCtrl := difficultyController.NewDifficultyNewsController(db)

	// 🔐 Role admin, teacher, owner
	protectedDifficultyRoutes := api.Group("/difficulties",
		authMiddleware.OnlyRoles(
			constants.RoleErrorTeacher("mengelola difficulties"),
			constants.TeacherAndAbove...,
		),
	)
	protectedDifficultyRoutes.Post("/", difficultyCtrl.CreateDifficulty)
	protectedDifficultyRoutes.Put("/:id", difficultyCtrl.UpdateDifficulty)
	protectedDifficultyRoutes.Delete("/:id", difficultyCtrl.DeleteDifficulty)

	protectedDifficultyNewsRoutes := api.Group("/difficulties-news",
		authMiddleware.OnlyRoles(
			constants.RoleErrorTeacher("mengelola difficulty news"),
			constants.TeacherAndAbove...,
		),
	)
	protectedDifficultyNewsRoutes.Post("/", difficultyNewsCtrl.CreateNews)
	protectedDifficultyNewsRoutes.Put("/:id", difficultyNewsCtrl.UpdateNews)
	protectedDifficultyNewsRoutes.Delete("/:id", difficultyNewsCtrl.DeleteNews)
}
