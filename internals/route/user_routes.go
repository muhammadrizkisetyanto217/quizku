package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	userRoute "quizku/internals/features/users/user/route"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
	userRoute.UserRoutes(app, db)
}
