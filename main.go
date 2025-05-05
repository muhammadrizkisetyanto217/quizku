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
	// ✅ Load .env variables
	configs.LoadEnv()

	// ✅ Inisialisasi Fiber
	app := fiber.New()

	// ✅ Setup global middleware (logger, recovery, dll)
	middlewares.SetupMiddlewares(app)

	// ✅ Koneksi ke database
	database.ConnectDB()

	// ✅ Jalankan scheduler pembersih token blacklist
	scheduler.StartBlacklistCleanupScheduler(database.DB)

	// ✅ Setup semua route
	routes.SetupRoutes(app, database.DB)

	// ✅ Ambil PORT dari Railway atau default 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// ✅ Jalankan aplikasi
	log.Printf("🚀 Server running at http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
