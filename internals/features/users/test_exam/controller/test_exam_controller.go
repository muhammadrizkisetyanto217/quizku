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

// ✅ Get all test exams
func (ctrl *TestExamController) GetAll(c *fiber.Ctx) error {
	var exams []model.TestExam
	if err := ctrl.DB.Order("id DESC").Find(&exams).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch test exams"})
	}
	return c.JSON(exams)
}

// ✅ Get test exam by ID
func (ctrl *TestExamController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var exam model.TestExam
	if err := ctrl.DB.First(&exam, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Test exam not found"})
	}
	return c.JSON(exam)
}

// ✅ Create new test exam
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

// ✅ Update test exam
func (ctrl *TestExamController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	// Cari data berdasarkan ID
	var exam model.TestExam
	if err := ctrl.DB.First(&exam, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Test exam not found"})
	}

	// Parse input dari user
	var payload model.TestExam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi status agar tidak kosong dan sesuai nilai yang diizinkan
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"archived": true,
	}

	if payload.Status == "" {
		payload.Status = "active" // fallback default
	} else if !validStatuses[payload.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status value"})
	}

	// Update field
	exam.Name = payload.Name
	exam.Status = payload.Status
	exam.UpdatedAt = time.Now()

	// Simpan ke DB
	if err := ctrl.DB.Save(&exam).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update test exam"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Test exam updated successfully",
		"data":    exam,
	})
}

// ✅ Delete test exam
func (ctrl *TestExamController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.TestExam{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete test exam"})
	}
	return c.JSON(fiber.Map{"message": "Test exam deleted successfully"})
}
