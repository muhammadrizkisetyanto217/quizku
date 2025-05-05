package controller

import (
	"quizku/internals/features/users/survey/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SurveyQuestionController struct {
	DB *gorm.DB
}

func NewSurveyQuestionController(db *gorm.DB) *SurveyQuestionController {
	return &SurveyQuestionController{DB: db}
}

// ✅ Get all questions
func (ctrl *SurveyQuestionController) GetAll(c *fiber.Ctx) error {
	var questions []model.SurveyQuestion
	if err := ctrl.DB.Order("order_index ASC").Find(&questions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch questions"})
	}
	return c.JSON(questions)
}

// ✅ Get question by ID
func (ctrl *SurveyQuestionController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var question model.SurveyQuestion
	if err := ctrl.DB.First(&question, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Question not found"})
	}
	return c.JSON(question)
}

// ✅ Create new question
func (ctrl *SurveyQuestionController) Create(c *fiber.Ctx) error {
	var payload model.SurveyQuestion
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Hitung order_index terakhir
	var maxOrder int
	ctrl.DB.Model(&model.SurveyQuestion{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)
	payload.OrderIndex = maxOrder + 1

	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create question"})
	}
	return c.Status(201).JSON(payload)
}

// ✅ Update question
func (ctrl *SurveyQuestionController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var question model.SurveyQuestion
	if err := ctrl.DB.First(&question, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Question not found"})
	}

	var payload model.SurveyQuestion
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	question.QuestionText = payload.QuestionText
	question.QuestionAnswer = payload.QuestionAnswer
	question.UpdatedAt = time.Now()

	if err := ctrl.DB.Save(&question).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update question"})
	}

	return c.JSON(question)
}

// ✅ Delete question
func (ctrl *SurveyQuestionController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.SurveyQuestion{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete question"})
	}
	return c.JSON(fiber.Map{"message": "Question deleted successfully"})
}
