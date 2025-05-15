package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"quizku/internals/features/certificates/certificate_versions/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CertificateVersionController struct {
	DB *gorm.DB
}

func NewCertificateVersionController(db *gorm.DB) *CertificateVersionController {
	return &CertificateVersionController{DB: db}
}

func (ctrl *CertificateVersionController) GetAll(c *fiber.Ctx) error {
	var versions []model.CertificateVersionModel
	if err := ctrl.DB.Find(&versions).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil data versi sertifikat"})
	}
	return c.JSON(fiber.Map{"data": versions})
}

func (ctrl *CertificateVersionController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var version model.CertificateVersionModel
	if err := ctrl.DB.First(&version, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Versi sertifikat tidak ditemukan"})
	}
	return c.JSON(version)
}

func (ctrl *CertificateVersionController) Create(c *fiber.Ctx) error {
	var payload model.CertificateVersionModel
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	payload.CreatedAt = time.Now()
	// payload.UpdatedAt = time.Now()

	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat versi sertifikat"})
	}
	return c.JSON(payload)
}

func (ctrl *CertificateVersionController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var version model.CertificateVersionModel
	if err := ctrl.DB.First(&version, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Versi tidak ditemukan"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	// Siapkan map kosong untuk field yang akan diupdate
	updates := map[string]interface{}{}

	if note, ok := updateData["note"].(string); ok {
		updates["note"] = note
	}

	if totalThemes, ok := updateData["total_themes"].(float64); ok {
		updates["total_themes"] = int(totalThemes) // karena JSON number = float64
	}

	now := time.Now()
	updates["updated_at"] = now

	if len(updates) == 1 { // hanya updated_at doang? berarti tidak ada perubahan signifikan
		return c.JSON(fiber.Map{"message": "Tidak ada perubahan data"})
	}

	if err := ctrl.DB.Model(&version).Updates(updates).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update versi"})
	}

	return c.JSON(version)
}

func (ctrl *CertificateVersionController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var version model.CertificateVersionModel
	if err := ctrl.DB.First(&version, idInt).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Versi sertifikat dengan ID %s tidak ditemukan", id)})
	}

	if err := ctrl.DB.Delete(&version).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Gagal menghapus versi dengan ID %s", id)})
	}

	return c.JSON(fiber.Map{
		"message":    fmt.Sprintf("Versi sertifikat dengan ID %s berhasil dihapus", id),
		"deleted_id": id,
	})
}
