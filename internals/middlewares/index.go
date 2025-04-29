package middlewares

import (
	loggerMiddleware "quizku/internals/middlewares/logger"

	"github.com/gofiber/fiber/v2"
)

// SetupMiddlewares menggabungkan semua middleware penting
func SetupMiddlewares(app *fiber.App) {
	app.Use(RecoveryMiddleware())                // 🔥 Tangkap panic
	app.Use(loggerMiddleware.LoggerMiddleware()) // 📝 Logger Request/Response
	app.Use(CorsMiddleware())                    // 🌐 CORS global
}
