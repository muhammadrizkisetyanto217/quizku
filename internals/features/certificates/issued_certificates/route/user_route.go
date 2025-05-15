// routes/issued_certificate_user_routes.go
package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	controller "quizku/internals/features/certificates/issued_certificates/controller"
)

func IssuedCertificateUserRoutes(app fiber.Router, db *gorm.DB) {
	c := controller.NewIssuedCertificateController(db)
	app.Get("/me", c.GetByID) // GET all certificates for current user
	app.Get("/category/:subcategory_id", c.GetBySubcategoryID) // GET all certificates for current user by subcategory
}
