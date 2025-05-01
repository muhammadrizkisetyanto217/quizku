package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	authRoute "quizku/internals/features/users/auth/route"
)

func AuthRoutes(app *fiber.App, db *gorm.DB) {

	authRoute.AuthRoutes(app, db)

}
