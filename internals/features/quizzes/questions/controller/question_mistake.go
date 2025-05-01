package controller

import (
	"log"
	"time"

	questionMistakeModel "quizku/internals/features/quizzes/questions/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuestionMistakeController struct {
	DB *gorm.DB
}

func NewQuestionMistakeController(db *gorm.DB) *QuestionMistakeController {
	return &QuestionMistakeController{DB: db}
}

// ✅ CREATE /api/question-mistakes
func (ctrl *QuestionMistakeController) Create(c *fiber.Ctx) error {
	start := time.Now()
	log.Println("[START] CreateQuestionMistake")

	var single questionMistakeModel.QuestionMistakeModel
	var multiple []questionMistakeModel.QuestionMistakeModel

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// Jika data berupa array
		if err := c.BodyParser(&multiple); err != nil {
			log.Println("[ERROR] Failed to parse array:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid array format"})
		}
		if len(multiple) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Array is empty"})
		}
		if err := ctrl.DB.Create(&multiple).Error; err != nil {
			log.Println("[ERROR] Failed to insert multiple question_mistakes:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Insert failed"})
		}
		log.Printf("[DONE] Created %d mistakes in %.2fms", len(multiple), time.Since(start).Seconds()*1000)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple question mistakes saved",
			"data":    multiple,
		})
	}

	// Jika data tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body format"})
	}
	if err := ctrl.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to insert question_mistake:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Insert failed"})
	}

	log.Printf("[DONE] Created mistake in %.2fms", time.Since(start).Seconds()*1000)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Question mistake saved",
		"data":    single,
	})
}

// ✅ GET /api/question-mistakes/:user_id
func (ctrl *QuestionMistakeController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var mistakes []questionMistakeModel.QuestionMistakeModel

	if err := ctrl.DB.Where("user_id = ?", userID).Find(&mistakes).Error; err != nil {
		log.Println("[ERROR] Failed to get mistakes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve question mistakes"})
	}

	return c.JSON(mistakes)
}

// ✅ DELETE /api/question-mistakes/:id
func (ctrl *QuestionMistakeController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var mistake questionMistakeModel.QuestionMistakeModel

	if err := ctrl.DB.First(&mistake, id).Error; err != nil {
		log.Println("[ERROR] Question mistake not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question mistake not found"})
	}

	if err := ctrl.DB.Delete(&mistake).Error; err != nil {
		log.Println("[ERROR] Failed to delete question mistake:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete question mistake"})
	}

	return c.JSON(fiber.Map{"message": "Question mistake deleted successfully"})
}
