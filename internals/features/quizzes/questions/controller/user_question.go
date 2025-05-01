package controller


import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	userQuestionModel "quizku/internals/features/quizzes/questions/model"

)

type UserQuestionController struct {
	DB *gorm.DB
}

func NewUserQuestionController(db *gorm.DB) *UserQuestionController {
	return &UserQuestionController{DB: db}
}

// POST /api/user_questions
func (ctrl *UserQuestionController) Create(c *fiber.Ctx) error {
	start := time.Now()
	log.Println("[START] Create UserQuestion")

	var single userQuestionModel.UserQuestionModel
	var multiple []userQuestionModel.UserQuestionModel

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		if err := c.BodyParser(&multiple); err != nil {
			log.Println("[ERROR] Failed to parse array:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid array format",
			})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Array is empty")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Request array is empty",
			})
		}

		if err := ctrl.DB.Create(&multiple).Error; err != nil {
			log.Println("[ERROR] Failed to insert multiple user_questions:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Insert failed",
			})
		}

		log.Printf("[DONE] %d UserQuestions created in %.2fms", len(multiple), time.Since(start).Seconds()*1000)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple user questions saved",
			"data":    multiple,
		})
	}

	// Fallback: single object
	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single user_question:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid body format",
		})
	}

	if err := ctrl.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to insert single user_question:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Insert failed",
		})
	}

	log.Printf("[DONE] Single UserQuestion created in %.2fms", time.Since(start).Seconds()*1000)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User question saved",
		"data":    single,
	})
}


// GET /api/user_questions/user/:user_id
func (ctrl *UserQuestionController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var results []userQuestionModel.UserQuestionModel

	if err := ctrl.DB.Where("user_id = ?", userID).Find(&results).Error; err != nil {
		log.Println("[ERROR] Failed to fetch user_questions by user_id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user questions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User questions fetched successfully",
		"data":    results,
	})
}


// GET /api/user_questions/user/:user_id/question/:question_id
func (ctrl *UserQuestionController) GetByUserIDAndQuestionID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	questionID := c.Params("question_id")

	var result userQuestionModel.UserQuestionModel

	if err := ctrl.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&result).Error; err != nil {
		log.Println("[ERROR] User question not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User question not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User question fetched successfully",
		"data":    result,
	})
}
