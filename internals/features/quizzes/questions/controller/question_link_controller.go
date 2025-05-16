package controller

import (
	"fmt"
	"log"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	link := model.QuestionLink{
		QuestionID: req.QuestionID,
		TargetType: req.TargetType, // e.g. "quiz", "exam", "reading"
		TargetID:   req.TargetID,
	}

	if err := ctrl.DB.Create(&link).Error; err != nil {
		log.Println("[ERROR] Gagal membuat question link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}

	return c.JSON(fiber.Map{"message": "Link berhasil dibuat", "data": link})
}

// 🟢 GET /api/question-links
// Mengambil semua data `question_link` tanpa filter.
// Cocok untuk admin panel atau validasi hubungan soal secara global.
func (ctrl *QuestionLinkController) GetAll(c *fiber.Ctx) error {
	var links []model.QuestionLink
	if err := ctrl.DB.Find(&links).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	return c.JSON(fiber.Map{"total": len(links), "data": links})
}

// 🟢 GET /api/question-links/question/:id
// Mengambil semua link soal berdasarkan question_id.
// Cocok untuk menampilkan hubungan/histori dari satu soal tertentu.
func (ctrl *QuestionLinkController) GetByQuestionID(c *fiber.Ctx) error {
	questionID := c.Params("id")
	var links []model.QuestionLink
	if err := ctrl.DB.Where("question_id = ?", questionID).Find(&links).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	return c.JSON(links)
}

// 🟠 PUT /api/question-links/:id
// Mengupdate data link soal berdasarkan ID.
// Umumnya dipakai untuk memperbaiki target dari soal yang sudah ada.
func (ctrl *QuestionLinkController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.QuestionLinkRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var link model.QuestionLink
	if err := ctrl.DB.First(&link, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	link.QuestionID = req.QuestionID
	link.TargetType = req.TargetType
	link.TargetID = req.TargetID

	if err := ctrl.DB.Save(&link).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update data"})
	}

	return c.JSON(fiber.Map{"message": "Berhasil update", "data": link})
}

// 🔴 DELETE /api/question-links/:id
// Menghapus satu link soal berdasarkan ID.
// Hati-hati karena ini akan memutus keterkaitan antara soal dengan entitas target.
func (ctrl *QuestionLinkController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var link model.QuestionLink

	if err := ctrl.DB.First(&link, id).Error; err != nil {
		log.Println("[ERROR] Question link not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Question link not found",
		})
	}

	if err := ctrl.DB.Delete(&link).Error; err != nil {
		log.Println("[ERROR] Failed to delete question link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete question link",
		})
	}

	log.Printf("[SUCCESS] Question link with ID %v deleted\n", link.ID)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Question link with ID %v deleted successfully", link.ID),
	})
}
