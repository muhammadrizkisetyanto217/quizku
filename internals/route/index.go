package routes

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	routeDetails "quizku/internals/route/details"
)

var startTime time.Time

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	startTime = time.Now()

	BaseRoutes(app, db)

	log.Println("[INFO] Setting up AuthRoutes...")
	routeDetails.AuthRoutes(app, db)

	log.Println("[INFO] Setting up UserRoutes...")
	routeDetails.UserRoutes(app, db)

	log.Println("[INFO] Setting up LessonRoutes...")
	routeDetails.LessonRoutes(app, db)

	log.Println("[INFO] Setting up QuizzesRoute...")
	routeDetails.QuizzesRoute(app, db)
} 