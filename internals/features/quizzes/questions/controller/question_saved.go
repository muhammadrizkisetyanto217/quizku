package controller

import (
	"fmt"
	"log"

	questionModel "quizku/internals/features/quizzes/questions/model"
	questionSavedModel "quizku/internals/features/quizzes/questions/model"

	"github.com/gofiber/fiber/v2"
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
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Saved multiple questions", "data": multiple})
	}

	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body format"})
	}
	if err := ctrl.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to insert question_saved:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Insert failed"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Question saved", "data": single})
}

// 🔹 GET /api/question-saved/:user_id
// Mengambil semua soal yang disimpan user berdasarkan user_id.
// Cocok untuk halaman "Soal Favorit Saya".
func (ctrl *QuestionSavedController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	log.Printf("[INFO] Fetching question_saved for user: %s", userID)

	var saved []questionSavedModel.QuestionSavedModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&saved).Error; err != nil {
		log.Println("[ERROR] Failed to fetch:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	return c.JSON(saved)
}

// 🔹 GET /api/question-saved/:user_id/full
// Mengambil daftar soal yang disimpan user, lengkap dengan data soalnya.
// Cocok untuk frontend yang ingin langsung menampilkan detail soalnya juga.
func (ctrl *QuestionSavedController) GetByUserIDWithQuestions(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	log.Printf("[INFO] Fetching question_saved WITH questions for user: %s", userID)

	var saved []questionSavedModel.QuestionSavedModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&saved).Error; err != nil {
		log.Println("[ERROR] Failed to fetch question_saved:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch question_saved"})
	}

	// Ambil daftar ID soal dari data saved
	var questionIDs []uint
	for _, s := range saved {
		questionIDs = append(questionIDs, s.QuestionID)
	}

	// Ambil data detail soalnya
	var questions []questionModel.QuestionModel
	if err := ctrl.DB.Where("id IN ?", questionIDs).Find(&questions).Error; err != nil {
		log.Println("[ERROR] Failed to fetch questions:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch questions"})
	}

	// Gabungkan data saved + soal
	type Combined struct {
		questionSavedModel.QuestionSavedModel
		Question questionModel.QuestionModel `json:"question"`
	}

	var combined []Combined
	questionMap := map[uint]questionModel.QuestionModel{}
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for _, s := range saved {
		if question, ok := questionMap[s.QuestionID]; ok {
			combined = append(combined, Combined{
				QuestionSavedModel: s,
				Question:           question,
			})
		}
	}

	return c.JSON(combined)
}

// 🔹 DELETE /api/question-saved/:id
// Menghapus satu data soal yang disimpan berdasarkan ID.
// Cocok digunakan saat user ingin menghapus soal dari daftar favorit.
func (ctrl *QuestionSavedController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting question_saved with ID: %s", id)

	if err := ctrl.DB.Delete(&questionSavedModel.QuestionSavedModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete",
		})
	}

	log.Printf("[SUCCESS] question_saved with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("question_saved with ID %s deleted successfully", id),
	})
}
