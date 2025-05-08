package details

import (
	donationController "quizku/internals/features/donations/donations/controller"
	donationRoutes "quizku/internals/features/donations/donations/routes"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DonationRoutes(app *fiber.App, db *gorm.DB) {
	// 🔐 Semua route donasi yang aman
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 👤 Prefix user
	userGroup := api.Group("/u")
	donationRoutes.DonationRoutes(userGroup, db)

	// 🔓 Webhook Midtrans tanpa middleware (public)
	app.Post("/api/donations/notification", func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return donationController.NewDonationController(db).HandleMidtransNotification(c)
	})
}
