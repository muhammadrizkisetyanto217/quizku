package route

import (
	readingController "quizku/internals/features/quizzes/readings/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReadingUserRoutes(api fiber.Router, db *gorm.DB) {
	readingCtrl := readingController.NewReadingController(db)
	userReadingCtrl := readingController.NewUserReadingController(db)

	readingRoutes := api.Group("/readings")
	readingRoutes.Get("/", readingCtrl.GetReadings)
	readingRoutes.Get("/:id", readingCtrl.GetReading)
	readingRoutes.Get("/unit/:unitId", readingCtrl.GetReadingsByUnit)

	// Tooltips & Konversi
	readingRoutes.Get("/:id/with-tooltips", readingCtrl.GetReadingWithTooltips)
	

	userReadingRoutes := api.Group("/user-readings")
	userReadingRoutes.Post("/", userReadingCtrl.CreateUserReading)
	userReadingRoutes.Get("/:id", userReadingCtrl.GetAllUserReading)
	userReadingRoutes.Get("/user/:user_id", userReadingCtrl.GetByUserID)
}