package details

import (
	certRoute "quizku/internals/features/certificates/certificate_versions/route"
	issuedCertRoute "quizku/internals/features/certificates/issued_certificates/route"
	rateLimiter "quizku/internals/middlewares"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CertificateRoutes(app *fiber.App, db *gorm.DB) {
	// 🔐 Semua route aman (butuh token + rate limit)
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 🔐 Admin routes
	adminGroup := api.Group("/a")
	certRoute.CertificateVersionAdminRoutes(adminGroup, db)
	issuedCertRoute.IssuedCertificateAdminRoutes(adminGroup.Group("/certificates"), db)

	// 👤 User routes
	userGroup := api.Group("/u")
	issuedCertRoute.IssuedCertificateUserRoutes(userGroup.Group("/certificates"), db)

	// 🔓 Public route (tidak pakai middleware)
	issuedCertRoute.IssuedCertificatePublicRoutes(app, db)
}
