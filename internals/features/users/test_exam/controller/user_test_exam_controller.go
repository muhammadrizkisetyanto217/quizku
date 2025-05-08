package controller

import (
	"quizku/internals/features/users/test_exam/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserTestExamController struct {
	DB *gorm.DB
}

func NewUserTestExamController(db *gorm.DB) *UserTestExamController {
	return &UserTestExamController{DB: db}
}

// ✅ Get all user test exam
func (ctrl *UserTestExamController) GetAll(c *fiber.Ctx) error {
	var results []model.UserTestExam
	if err := ctrl.DB.Order("id DESC").Find(&results).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user test exam data"})
	}
	return c.JSON(results)
}

// ✅ Get by ID
func (ctrl *UserTestExamController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var data model.UserTestExam
	if err := ctrl.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User test exam not found"})
	}
	return c.JSON(data)
}

// ✅ Create new
func (ctrl *UserTestExamController) Create(c *fiber.Ctx) error {
	var payload model.UserTestExam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	payload.CreatedAt = time.Now()

	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user test exam"})
	}
	return c.Status(201).JSON(payload)
}

// ✅ Get by user_id (misalnya untuk riwayat nilai user)
func (ctrl *UserTestExamController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var results []model.UserTestExam
	if err := ctrl.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&results).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user test exam by user ID"})
	}
	return c.JSON(results)
}

// ✅ Get by test_exam_id (misalnya untuk nilai seluruh peserta di 1 ujian)
func (ctrl *UserTestExamController) GetByTestExamID(c *fiber.Ctx) error {
	examID := c.Params("test_exam_id")
	var results []model.UserTestExam
	if err := ctrl.DB.Where("test_exam_id = ?", examID).Order("created_at DESC").Find(&results).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user test exam by exam ID"})
	}
	return c.JSON(results)
}
