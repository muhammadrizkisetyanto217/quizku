package details

import (
	donationController "quizku/internals/features/donations/donations/controller"
	donationRoutes "quizku/internals/features/donations/donations/routes"
	donationQuestionAdminRoutes "quizku/internals/features/donations/donation_questions/route"
	donationQuestionUserRoutes "quizku/internals/features/donations/donation_questions/route"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DonationRoutes(app *fiber.App, db *gorm.DB) {
	// Semua route aman → membutuhkan token + rate limit
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 👤 Route untuk user biasa (/api/u)
	userGroup := api.Group("/u")
	donationRoutes.DonationRoutes(userGroup, db) // data donasi user
	donationQuestionUserRoutes.DonationQuestionUserRoutes(userGroup.Group("/donation-questions"), db)

	// 🔐 Route untuk admin/owner (/api/a)
	adminGroup := api.Group("/a")
	donationQuestionAdminRoutes.DonationQuestionAdminRoutes(adminGroup.Group("/donation-questions"), db)

	// 🔓 Webhook dari Midtrans (tidak pakai middleware)
	app.Post("/api/donations/notification", func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return donationController.NewDonationController(db).HandleMidtransNotification(c)
	})
}
