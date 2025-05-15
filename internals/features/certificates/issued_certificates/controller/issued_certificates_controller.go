package controller

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	model "quizku/internals/features/certificates/issued_certificates/model"
)

type IssuedCertificateController struct {
	DB *gorm.DB
}

func NewIssuedCertificateController(db *gorm.DB) *IssuedCertificateController {
	return &IssuedCertificateController{DB: db}
}

// ✅ GET /api/certificates/:id
func (ctrl *IssuedCertificateController) GetByIDUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	var cert model.IssuedCertificateModel
	if err := ctrl.DB.First(&cert, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Sertifikat tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Detail sertifikat ditemukan",
		"data":    cert,
	})
}

// ✅ Untuk User: Get all certificates miliknya sendiri
func (ctrl *IssuedCertificateController) GetByID(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id format"})
	}

	var certificates []model.IssuedCertificateModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&certificates).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil data sertifikat"})
	}

	return c.JSON(fiber.Map{"data": certificates})
}

// ✅ Untuk Public: Get certificate by slug (tanpa login)
func (ctrl *IssuedCertificateController) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Slug tidak boleh kosong"})
	}

	var cert model.IssuedCertificateModel
	if err := ctrl.DB.Where("slug_url = ?", slug).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sertifikat tidak ditemukan"})
	}

	return c.JSON(cert)
}
