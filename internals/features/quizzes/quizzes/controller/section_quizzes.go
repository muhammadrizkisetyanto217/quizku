package controller

import (
	"fmt"
	"log"

	"quizku/internals/features/quizzes/quizzes/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SectionQuizController struct {
	DB *gorm.DB
}

func NewSectionQuizController(db *gorm.DB) *SectionQuizController {
	return &SectionQuizController{DB: db}
}

// ✅ GET /api/section-quizzes
// Mengambil seluruh data section_quizzes dari database.
func (sqc *SectionQuizController) GetSectionQuizzes(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all section quizzes")

	var quizzes []model.SectionQuizzesModel
	if err := sqc.DB.Find(&quizzes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch section quizzes:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch section quizzes",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d section quizzes\n", len(quizzes))
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Section quizzes fetched successfully",
		"total":   len(quizzes),
		"data":    quizzes,
	})
}

// ✅ GET /api/section-quizzes/:id
// Mengambil satu data section_quiz berdasarkan ID-nya.

func (sqc *SectionQuizController) GetSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching section quiz with ID: %s\n", id)

	var quiz model.SectionQuizzesModel
	if err := sqc.DB.
		Where("section_quizzes_id = ?", id).
		First(&quiz).Error; err != nil {
		log.Println("[ERROR] Section quiz not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Section quiz not found",
		})
	}

	log.Printf("[SUCCESS] Retrieved section quiz with ID: %s\n", id)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Section quiz fetched successfully",
		"data":    quiz,
	})
}

// ✅ GET /api/section-quizzes/unit/:unitId
// Mengambil semua section_quizzes yang termasuk dalam unit tertentu.
func (sqc *SectionQuizController) GetSectionQuizzesByUnit(c *fiber.Ctx) error {
	unitID := c.Params("unitId")
	log.Printf("[INFO] Fetching section quizzes for section_quizzes_unit_id: %s\n", unitID)

	var sectionQuizzes []model.SectionQuizzesModel
	if err := sqc.DB.
		Where("section_quizzes_unit_id = ?", unitID).
		Find(&sectionQuizzes).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch section quizzes for unit_id %s: %v\n", unitID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch section quizzes by unit ID",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d section quizzes for section_quizzes_unit_id %s\n", len(sectionQuizzes), unitID)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Section quizzes fetched by unit ID successfully",
		"total":   len(sectionQuizzes),
		"data":    sectionQuizzes,
	})
}

// ✅ POST /api/section-quizzes
// Membuat satu atau banyak data section_quizzes sekaligus.
func (sqc *SectionQuizController) CreateSectionQuiz(c *fiber.Ctx) error {
	log.Println("[INFO] Creating section quiz (single or multiple)")

	var single model.SectionQuizzesModel
	var multiple []model.SectionQuizzesModel

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// 🟡 JSON berupa array → batch insert
		if err := c.BodyParser(&multiple); err != nil {
			log.Println("[ERROR] Failed to parse section quizzes array:", err)
			return c.Status(400).JSON(fiber.Map{"status": false, "message": "Invalid array request"})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Received empty array")
			return c.Status(400).JSON(fiber.Map{"status": false, "message": "Request array is empty"})
		}

		// 🛠️ Tambahan validasi semantik (opsional)

		if err := sqc.DB.Create(&multiple).Error; err != nil {
			log.Println("[ERROR] Failed to create multiple section quizzes:", err)
			return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to create section quizzes"})
		}

		log.Printf("[SUCCESS] %d section quizzes created\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"status":  true,
			"message": "Section quizzes created successfully",
			"data":    multiple,
		})
	}

	// 🟢 Input berupa satu objek
	if err := c.BodyParser(&single); err != nil {
		log.Println("[ERROR] Failed to parse single section quiz:", err)
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Invalid request format (expected object or array)"})
	}

	// 🛠️ Validasi semantik (opsional)

	if err := sqc.DB.Create(&single).Error; err != nil {
		log.Println("[ERROR] Failed to create section quiz:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to create section quiz"})
	}

	log.Printf("[SUCCESS] Section quiz created with ID: %d\n", single.SectionQuizzesID)
	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "Section quiz created successfully",
		"data":    single,
	})
}

// ✅ PUT /api/section-quizzes/:id
// Melakukan update terhadap satu section_quiz berdasarkan ID.
func (sqc *SectionQuizController) UpdateSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating section quiz with ID: %s\n", id)

	var quiz model.SectionQuizzesModel
	if err := sqc.DB.First(&quiz, "section_quizzes_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Section quiz not found:", err)
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Section quiz not found"})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Invalid request"})
	}

	if err := sqc.DB.Model(&quiz).Updates(requestData).Error; err != nil {
		log.Println("[ERROR] Failed to update section quiz:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to update section quiz"})
	}

	log.Printf("[SUCCESS] Section quiz with ID %s updated\n", id)
	return c.JSON(fiber.Map{"status": true, "message": "Section quiz updated successfully", "data": quiz})
}

// ✅ DELETE /api/section-quizzes/:id
// Menghapus satu section_quiz berdasarkan ID.
func (sqc *SectionQuizController) DeleteSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting section quiz with ID: %s\n", id)

	if err := sqc.DB.Where("section_quizzes_id = ?", id).Delete(&model.SectionQuizzesModel{}).Error; err != nil {
		log.Println("[ERROR] Failed to delete section quiz:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to delete section quiz",
		})
	}

	log.Printf("[SUCCESS] Section quiz with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("Section quiz with ID %s deleted successfully", id),
	})
}
