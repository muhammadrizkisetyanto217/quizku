package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"quizku/internals/configs"
	database "quizku/internals/databases"
	"quizku/internals/features/donations/donations/service"
	scheduler "quizku/internals/features/users/auth/scheduler"
	middlewares "quizku/internals/middlewares"
	routes "quizku/internals/route"
)

func main() {
	configs.LoadEnv()
	app := fiber.New()

	middlewares.SetupMiddlewares(app)

	// ✅ Koneksi DB
	database.ConnectDB()
	scheduler.StartBlacklistCleanupScheduler(database.DB)

	// ✅ Ambil MIDTRANS_SERVER_KEY dari .env
	serverKey := configs.GetEnv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("❌ MIDTRANS_SERVER_KEY tidak ditemukan di .env")
	}

	service.InitMidtrans(serverKey) // ✅ PASANG PARAMETERNYA

	// ✅ Setup routes dulu
	routes.SetupRoutes(app, database.DB)

	// ✅ Baru tangani preflight OPTIONS (setelah semua route aktif)
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent) // 204 No Content
	})

	// ✅ Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("✅ Listening on PORT: %s", port)
	log.Fatal(app.Listen(":" + port))
}
