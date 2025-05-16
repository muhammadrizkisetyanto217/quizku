package controller

import (
	"quizku/internals/features/users/test_exam/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TestExamController struct {
	DB *gorm.DB
}

func NewTestExamController(db *gorm.DB) *TestExamController {
	return &TestExamController{DB: db}
}

// GetAll mengambil semua data test exam yang tersedia.
func (ctrl *TestExamController) GetAll(c *fiber.Ctx) error {
	var exams []model.TestExam
	if err := ctrl.DB.Order("test_exam_id DESC").Find(&exams).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch test exams"})
	}
	return c.JSON(exams)
}

// ✅ Ambil test_exam berdasarkan ID
func (ctrl *TestExamController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var exam model.TestExam
	if err := ctrl.DB.First(&exam, "test_exam_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Test exam tidak ditemukan",
		})
	}
	return c.JSON(exam)
}

// Create membuat entri test exam baru di database.
func (ctrl *TestExamController) Create(c *fiber.Ctx) error {
	var payload model.TestExam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create test exam"})
	}
	return c.Status(201).JSON(payload)
}

// ✅ Perbarui data test_exam berdasarkan ID
func (ctrl *TestExamController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var exam model.TestExam
	if err := ctrl.DB.First(&exam, "test_exam_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Test exam tidak ditemukan",
		})
	}

	var payload model.TestExam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Body request tidak valid"})
	}

	// Validasi status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"archived": true,
	}
	if payload.TestExamStatus != "" && !validStatuses[payload.TestExamStatus] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status tidak valid (harus: active, pending, archived)",
		})
	}

	exam.TestExamName = payload.TestExamName
	if payload.TestExamStatus != "" {
		exam.TestExamStatus = payload.TestExamStatus
	}
	exam.UpdatedAt = time.Now()

	if err := ctrl.DB.Save(&exam).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui test_exam",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Test exam berhasil diperbarui",
		"data":    exam,
	})
}

// ✅ Hapus test_exam berdasarkan ID
func (ctrl *TestExamController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.TestExam{}, "test_exam_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus test_exam",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Test exam berhasil dihapus",
	})
}
