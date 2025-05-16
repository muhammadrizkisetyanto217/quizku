// routes/issued_certificate_public_routes.go
package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	controller "quizku/internals/features/certificates/user_certificates/controller"
)

func IssuedCertificatePublicRoutes(app *fiber.App, db *gorm.DB) {
	c := controller.NewIssuedCertificateController(db)
	app.Get("/certificates/:slug", c.GetBySlug) // public access via certificate_slug_url
}
