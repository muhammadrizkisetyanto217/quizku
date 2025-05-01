package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	difficultiesRoute "quizku/internals/features/lessons/difficulty/route"
)

func LessonRoutes(app *fiber.App, db *gorm.DB) {
	difficultiesRoute.CategoryRoutes(app, db)
}
