package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	questionMistakeModel "quizku/internals/features/quizzes/questions/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionMistakeController struct {
	DB *gorm.DB
}

func NewQuestionMistakeController(db *gorm.DB) *QuestionMistakeController {
	return &QuestionMistakeController{DB: db}
}

// ✅ POST /api/question-mistakes
// Menyimpan satu atau banyak kesalahan jawaban user terhadap soal.
// Bisa digunakan untuk tracking soal yang sering salah dijawab.
// Mendukung format tunggal dan array secara otomatis.
func (ctrl *QuestionMistakeController) Create(c *fiber.Ctx) error {
	start := time.Now()
	log.Println("[START] CreateQuestionMistake")

	var single questionMistakeModel.QuestionMistakeModel
	var multiple []questionMistakeModel.QuestionMistakeModel

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// 📦 Input berupa array
		if err := c.BodyParser(&multiple); err != nil {
			log.Println("[ERROR] Failed to parse array:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid array format",
			})
		}
		if len(multiple) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Array is empty",
			})
		}
		if err := ctrl.DB.Create(&multiple).Error; err != nil {
			log.Println("[ERROR] Failed to insert multiple question_mistakes:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Insert failed",
			})
		}

		log.Printf("[DONE] Created %d mistakes in %.2fms", len(multiple), time.Since(start).Seconds()*1000)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple question mistakes saved",
			"total":   len(multiple),
			"data":    multiple,
		})
	}

	// 🧾 Input berupa satu object
	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid body format",
		})
	}

	if err := ctrl.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to insert question_mistake:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Insert failed",
		})
	}

	log.Printf("[DONE] Created mistake in %.2fms", time.Since(start).Seconds()*1000)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Question mistake saved",
		"data":    single,
	})
}

// ✅ GET /api/question-mistakes/:user_id
// Mengambil semua kesalahan jawaban berdasarkan user_id.
// Berguna untuk membuat fitur "soal yang sering kamu salah jawab".
func (ctrl *QuestionMistakeController) GetByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid user_id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id format",
		})
	}

	var mistakes []questionMistakeModel.QuestionMistakeModel
	if err := ctrl.DB.
		Where("question_mistake_user_id = ?", userID).
		Find(&mistakes).Error; err != nil {
		log.Println("[ERROR] Failed to get question mistakes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve question mistakes",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d question mistakes for user %s", len(mistakes), userIDStr)
	return c.JSON(fiber.Map{
		"total": len(mistakes),
		"data":  mistakes,
	})
}

// ✅ DELETE /api/question-mistakes/:id
// Menghapus satu kesalahan soal berdasarkan ID.
// Umumnya digunakan oleh admin untuk bersih-bersih data kesalahan yang salah input.
func (ctrl *QuestionMistakeController) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[ERROR] Invalid ID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var mistake questionMistakeModel.QuestionMistakeModel
	if err := ctrl.DB.First(&mistake, "question_mistake_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Question mistake not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Question mistake not found",
		})
	}

	if err := ctrl.DB.Delete(&mistake).Error; err != nil {
		log.Println("[ERROR] Failed to delete question mistake:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete question mistake",
		})
	}

	log.Printf("[SUCCESS] Question mistake with ID %d deleted", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Question mistake with ID %d deleted successfully", id),
	})
}
