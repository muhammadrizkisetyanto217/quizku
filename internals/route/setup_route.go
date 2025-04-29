package routes

import (
	"log"
	"os"
	database "quizku/internals/databases"
	authRoute "quizku/internals/features/users/auth/route"
	userRoute "quizku/internals/features/users/user/route"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var startTime time.Time

// Register routes
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fiber & Supabase PostgreSQL connected successfully Sekarang 🚀")
	})

	app.Get("/panic-test", func(c *fiber.Ctx) error {
		panic("Simulasi panic error!") // sengaja panic
	})

	// ✨ Advanced Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := database.DB.DB()
		dbStatus := "Connected"
		serverStatus := "OK"
		httpStatus := fiber.StatusOK

		// 🚨 Kalau error DB, ubah status server
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "Database connection error"
			serverStatus = "DOWN"
			httpStatus = fiber.StatusServiceUnavailable // 503
		}

		uptime := time.Since(startTime).Seconds()

		return c.Status(httpStatus).JSON(fiber.Map{
			"status":         serverStatus,
			"database":       dbStatus,
			"server_time":    time.Now().Format(time.RFC3339),
			"uptime_seconds": int(uptime),
			"environment":    os.Getenv("RAILWAY_ENVIRONMENT"),
		})
	})

	log.Println("[INFO] Setting up AuthRoutes...")
	authRoute.AuthRoutes(app, db)

	log.Println("[INFO] Setting up UserRoutes...")
	userRoute.UserRoutes(app, db)

}
