package controller

import (
	"fmt"
	"log"
	"strconv"

	"quizku/internals/features/quizzes/questions/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuestionLinkController struct {
	DB *gorm.DB
}

func NewQuestionLinkController(db *gorm.DB) *QuestionLinkController {
	return &QuestionLinkController{DB: db}
}

// 🟡 POST /api/question-links
// Membuat satu data `question_link` baru yang menghubungkan soal ke entitas lain
// seperti quiz, exam, reading, evaluation, dsb. Berguna untuk strukturisasi soal.
func (ctrl *QuestionLinkController) Create(c *fiber.Ctx) error {
	var req model.QuestionLinkRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("[ERROR] Invalid request body for creating question link:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	link := model.QuestionLink{
		QuestionLinkQuestionID: req.QuestionLinkQuestionID,
		QuestionLinkTargetType: req.QuestionLinkTargetType,
		QuestionLinkTargetID:   req.QuestionLinkTargetID,
	}

	if err := ctrl.DB.Create(&link).Error; err != nil {
		log.Println("[ERROR] Failed to create question link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save question link",
		})
	}

	log.Printf("[SUCCESS] Question link created: %+v\n", link)
	return c.Status(201).JSON(fiber.Map{
		"message": "Question link created successfully",
		"data":    link,
	})
}

// 🟢 GET /api/question-links
// Mengambil semua data `question_link` tanpa filter.
// Cocok untuk admin panel atau validasi hubungan soal secara global.
func (ctrl *QuestionLinkController) GetAll(c *fiber.Ctx) error {
	var links []model.QuestionLink
	if err := ctrl.DB.Find(&links).Error; err != nil {
		log.Println("[ERROR] Failed to fetch all question links:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data",
		})
	}
	return c.JSON(fiber.Map{
		"total": len(links),
		"data":  links,
	})
}

// 🟢 GET /api/question-links/question/:id
// Mengambil semua link soal berdasarkan question_id.
// Cocok untuk menampilkan hubungan/histori dari satu soal tertentu.
func (ctrl *QuestionLinkController) GetByQuestionID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	questionID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("[ERROR] Invalid question_id parameter:", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid question_id",
		})
	}

	var links []model.QuestionLink
	if err := ctrl.DB.
		Where("question_link_question_id = ?", questionID).
		Find(&links).Error; err != nil {
		log.Println("[ERROR] Failed to fetch links by question_id:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(links),
		"data":  links,
	})
}

// 🟠 PUT /api/question-links/:id
// Mengupdate data link soal berdasarkan ID.
// Umumnya dipakai untuk memperbaiki target dari soal yang sudah ada.
func (ctrl *QuestionLinkController) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("[ERROR] Invalid ID parameter:", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req model.QuestionLinkRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var link model.QuestionLink
	if err := ctrl.DB.First(&link, "question_link_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Question link not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Question link not found"})
	}

	link.QuestionLinkQuestionID = req.QuestionLinkQuestionID
	link.QuestionLinkTargetType = req.QuestionLinkTargetType
	link.QuestionLinkTargetID = req.QuestionLinkTargetID

	if err := ctrl.DB.Save(&link).Error; err != nil {
		log.Println("[ERROR] Failed to update question link:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update data"})
	}

	log.Printf("[SUCCESS] Question link with ID %d updated\n", id)
	return c.JSON(fiber.Map{"message": "Successfully updated", "data": link})
}

// 🔴 DELETE /api/question-links/:id
// Menghapus satu link soal berdasarkan ID.
// Hati-hati karena ini akan memutus keterkaitan antara soal dengan entitas target.
func (ctrl *QuestionLinkController) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("[ERROR] Invalid ID parameter:", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var link model.QuestionLink
	if err := ctrl.DB.First(&link, "question_link_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Question link not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question link not found"})
	}

	if err := ctrl.DB.Delete(&link).Error; err != nil {
		log.Println("[ERROR] Failed to delete question link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete question link"})
	}

	log.Printf("[SUCCESS] Question link with ID %d deleted\n", id)
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Question link with ID %d deleted successfully", id)})
}
