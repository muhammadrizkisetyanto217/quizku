package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"quizku/internals/configs"
	database "quizku/internals/databases"
	scheduler "quizku/internals/features/users/auth/scheduler"
	middlewares "quizku/internals/middlewares"
	routes "quizku/internals/route"
)

func main() {
	configs.LoadEnv()
	app := fiber.New()

	// ✅ Aktifkan middleware lebih dulu
	middlewares.SetupMiddlewares(app)

	// ✅ Tangani preflight semua route sebelum SetupRoutes
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent) // 204 No Content
	})

	// ✅ Koneksi DB
	database.ConnectDB()
	scheduler.StartBlacklistCleanupScheduler(database.DB)

	// ✅ Route
	routes.SetupRoutes(app, database.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
