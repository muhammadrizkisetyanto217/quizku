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

	log.Println("[INFO] Setting up DonationRoutes...")
	routeDetails.DonationRoutes(app, db)

	log.Println("[INFO] Setting up UtilsRoutes...")
	routeDetails.UtilsRoutes(app, db)

	log.Println("[INFO] Setting up CertificateRoutes...")
	routeDetails.CertificateRoutes(app, db)

	log.Println("[INFO] Setting up ProgressRoutes...")
	routeDetails.ProgressRoutes(app, db)
} 