package routes

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var startTime time.Time

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	startTime = time.Now()

	BaseRoutes(app, db)

	log.Println("[INFO] Setting up AuthRoutes...")
	AuthRoutes(app, db)

	log.Println("[INFO] Setting up UserRoutes...")
	UserRoutes(app, db)

	log.Println("[INFO] Setting up LessonRoutes...")
	LessonRoutes(app, db)
}
