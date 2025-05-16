package controller

import (
	"fmt"
	"log"
	"strconv"

	questionModel "quizku/internals/features/quizzes/questions/model"
	questionSavedModel "quizku/internals/features/quizzes/questions/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionSavedController struct {
	DB *gorm.DB
}

func NewQuestionSavedController(db *gorm.DB) *QuestionSavedController {
	return &QuestionSavedController{DB: db}
}

// 🔹 POST /api/question-saved
// Menyimpan satu atau banyak soal ke daftar soal favorit (saved questions).
// Digunakan saat user ingin menyimpan soal tertentu untuk dipelajari ulang.
//
// ✅ Bisa input satu atau array langsung.
// ✅ Berguna untuk fitur "bookmark soal" di frontend.
func (ctrl *QuestionSavedController) Create(c *fiber.Ctx) error {
	log.Println("[INFO] Create QuestionSaved called")

	var single questionSavedModel.QuestionSavedModel
	var multiple []questionSavedModel.QuestionSavedModel

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// Input berupa array
		if err := c.BodyParser(&multiple); err != nil {
			log.Println("[ERROR] Failed to parse array:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid array format"})
		}
		if len(multiple) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Array is empty"})
		}
		if err := ctrl.DB.Create(&multiple).Error; err != nil {
			log.Println("[ERROR] Failed to insert multiple question_saved:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Insert failed"})
		}
		log.Printf("[SUCCESS] Inserted %d question_saved records", len(multiple))
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Saved multiple questions",
			"total":   len(multiple),
			"data":    multiple,
		})
	}

	// Input tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body format"})
	}
	if err := ctrl.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to insert question_saved:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Insert failed"})
	}

	log.Printf("[SUCCESS] Inserted question_saved ID: %d", single.QuestionSavedID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Question saved",
		"data":    single,
	})
}

// 🔹 GET /api/question-saved/:user_id
// Mengambil semua soal yang disimpan user berdasarkan user_id.
// Cocok untuk halaman "Soal Favorit Saya".
func (ctrl *QuestionSavedController) GetByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	log.Printf("[INFO] Fetching question_saved for user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id format"})
	}

	var saved []questionSavedModel.QuestionSavedModel
	if err := ctrl.DB.
		Where("question_saved_user_id = ?", userID).
		Find(&saved).Error; err != nil {
		log.Println("[ERROR] Failed to fetch question_saved:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	log.Printf("[SUCCESS] Retrieved %d saved questions for user %s", len(saved), userIDStr)
	return c.JSON(fiber.Map{
		"total": len(saved),
		"data":  saved,
	})
}

// 🔹 GET /api/question-saved/:user_id/full
// Mengambil daftar soal yang disimpan user, lengkap dengan data soalnya.
// Cocok untuk frontend yang ingin langsung menampilkan detail soalnya juga.
func (ctrl *QuestionSavedController) GetByUserIDWithQuestions(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	log.Printf("[INFO] Fetching question_saved WITH questions for user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid user_id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id format"})
	}

	var saved []questionSavedModel.QuestionSavedModel
	if err := ctrl.DB.
		Where("question_saved_user_id = ?", userID).
		Find(&saved).Error; err != nil {
		log.Println("[ERROR] Failed to fetch question_saved:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch question_saved"})
	}

	// Ambil daftar question_saved_question_id
	var questionIDs []uint
	for _, s := range saved {
		questionIDs = append(questionIDs, s.QuestionSavedQuestionID)
	}

	// Ambil data soalnya
	var questions []questionModel.QuestionModel
	if err := ctrl.DB.
		Where("question_id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		log.Println("[ERROR] Failed to fetch questions:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch questions"})
	}

	// Gabungkan data
	type Combined struct {
		questionSavedModel.QuestionSavedModel
		Question questionModel.QuestionModel `json:"question"`
	}

	questionMap := make(map[uint]questionModel.QuestionModel)
	for _, q := range questions {
		questionMap[q.QuestionID] = q
	}

	var combined []Combined
	for _, s := range saved {
		if q, ok := questionMap[s.QuestionSavedQuestionID]; ok {
			combined = append(combined, Combined{
				QuestionSavedModel: s,
				Question:           q,
			})
		}
	}

	log.Printf("[SUCCESS] Fetched %d combined question_saved entries", len(combined))
	return c.JSON(fiber.Map{
		"total": len(combined),
		"data":  combined,
	})
}

// 🔹 DELETE /api/question-saved/:id
// Menghapus satu data soal yang disimpan berdasarkan ID.
// Cocok digunakan saat user ingin menghapus soal dari daftar favorit.
func (ctrl *QuestionSavedController) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[ERROR] Invalid ID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	log.Printf("[INFO] Deleting question_saved with ID: %d", id)

	if err := ctrl.DB.
		Delete(&questionSavedModel.QuestionSavedModel{}, "question_saved_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Failed to delete:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete",
		})
	}

	log.Printf("[SUCCESS] question_saved with ID %d deleted", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("question_saved with ID %d deleted successfully", id),
	})
}
