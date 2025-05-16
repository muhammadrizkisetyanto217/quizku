package controller

import (
	"encoding/json"
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

// ✅ GetAll mengembalikan seluruh pertanyaan survei yang ada di database,
// diurutkan berdasarkan `order_index` secara ascending.
func (ctrl *SurveyQuestionController) GetAll(c *fiber.Ctx) error {
	var questions []model.SurveyQuestion
	if err := ctrl.DB.Order("order_index ASC").Find(&questions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch questions"})
	}
	return c.JSON(questions)
}

// ✅ GetByID mengambil satu data pertanyaan survei berdasarkan ID yang diberikan.
func (ctrl *SurveyQuestionController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var question model.SurveyQuestion
	if err := ctrl.DB.First(&question, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Question not found"})
	}
	return c.JSON(question)
}


// ✅ Create menambahkan satu atau banyak pertanyaan survei baru ke dalam database.
func (ctrl *SurveyQuestionController) Create(c *fiber.Ctx) error {
	body := c.Body()

	// Cek apakah body diawali dengan [ (berarti array)
	if len(body) > 0 && body[0] == '[' {
		var payloads []model.SurveyQuestion
		if err := json.Unmarshal(body, &payloads); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid array request body"})
		}

		// Ambil order terakhir
		var maxOrder int
		ctrl.DB.Model(&model.SurveyQuestion{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)

		for i := range payloads {
			payloads[i].OrderIndex = maxOrder + i + 1
			payloads[i].CreatedAt = time.Now()
			payloads[i].UpdatedAt = time.Now()
		}

		if err := ctrl.DB.Create(&payloads).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert questions"})
		}
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple questions created",
			"data":    payloads,
		})
	} else {
		var payload model.SurveyQuestion
		if err := json.Unmarshal(body, &payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid object request body"})
		}

		// Hitung order terakhir
		var maxOrder int
		ctrl.DB.Model(&model.SurveyQuestion{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)

		payload.OrderIndex = maxOrder + 1
		payload.CreatedAt = time.Now()
		payload.UpdatedAt = time.Now()

		if err := ctrl.DB.Create(&payload).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert question"})
		}
		return c.Status(201).JSON(fiber.Map{
			"message": "Single question created",
			"data":    payload,
		})
	}
}

// ✅ Update mengubah isi pertanyaan survei berdasarkan ID yang diberikan.
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

// ✅ Delete menghapus pertanyaan survei berdasarkan ID yang diberikan.
func (ctrl *SurveyQuestionController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.SurveyQuestion{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete question"})
	}
	return c.JSON(fiber.Map{"message": "Question deleted successfully"})
}
