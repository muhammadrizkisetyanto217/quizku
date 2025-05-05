// middlewares/cors.go

package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupMiddlewareCors membuat middleware CORS
func CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://quizku-production.vercel.app", // sesuaikan
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})
}
