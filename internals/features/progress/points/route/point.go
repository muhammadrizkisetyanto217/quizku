package routes

import (
	pointController "quizku/internals/features/progress/points/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserPointRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	userPointLogController := pointController.NewUserPointLogController(db)
	userPointRoutes := api.Group("/user-point-logs")

	userPointRoutes.Post("/", userPointLogController.Create)
	userPointRoutes.Get("/:user_id", userPointLogController.GetByUserID)
}
