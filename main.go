package main

import (
	"log"
	"quizku/internals/configs"
	database "quizku/internals/databases"
	scheduler "quizku/internals/features/users/auth/scheduler"
	"quizku/internals/middlewares"
	routes "quizku/internals/route"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// ✅ Muat file .env dulu
	configs.LoadEnv()
	// Inisialisasi Fiber
	app := fiber.New()

	middlewares.SetupMiddlewares(app)

	// Koneksi ke Supabase
	database.ConnectDB()

	// ✅ Jalankan scheduler harian
	scheduler.StartBlacklistCleanupScheduler(database.DB)

	// ✅ Panggil semua route dari folder routes
	routes.SetupRoutes(app, database.DB)

	// Start server
	log.Fatal(app.Listen(":3000"))
}