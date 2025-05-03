package controller

import (
	"fmt"
	"log"

	examModel "quizku/internals/features/quizzes/exams/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExamController struct {
	DB *gorm.DB
}

func NewExamController(db *gorm.DB) *ExamController {
	return &ExamController{DB: db}
}

// GET all exams
func (ec *ExamController) GetExams(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all exams")
	var exams []examModel.ExamModel

	if err := ec.DB.Find(&exams).Error; err != nil {
		log.Println("[ERROR] Failed to fetch exams:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch exams"})
	}

	log.Printf("[SUCCESS] Retrieved %d exams\n", len(exams))
	return c.JSON(fiber.Map{
		"message": "Exams fetched successfully",
		"total":   len(exams),
		"data":    exams,
	})
}

// GET exam by ID
func (ec *ExamController) GetExam(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching exam with ID:", id)

	var exam examModel.ExamModel
	if err := ec.DB.First(&exam, id).Error; err != nil {
		log.Println("[ERROR] Exam not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Exam not found"})
	}

	return c.JSON(fiber.Map{
		"message": "Exam fetched successfully",
		"data":    exam,
	})
}

// GET exams by unit_id
func (ec *ExamController) GetExamsByUnitID(c *fiber.Ctx) error {
	unitID := c.Params("unitId")
	log.Printf("[INFO] Fetching exams for unit_id: %s\n", unitID)

	var exams []examModel.ExamModel
	if err := ec.DB.Where("unit_id = ?", unitID).Find(&exams).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch exams for unit_id %s: %v\n", unitID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch exams"})
	}

	log.Printf("[SUCCESS] Retrieved %d exams for unit_id %s\n", len(exams), unitID)
	return c.JSON(fiber.Map{
		"message": "Exams fetched successfully by unit ID",
		"total":   len(exams),
		"data":    exams,
	})
}

// POST create a new exam
func (ec *ExamController) CreateExam(c *fiber.Ctx) error {
	log.Println("[INFO] Creating a new exam")

	var exam examModel.ExamModel
	if err := c.BodyParser(&exam); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := ec.DB.Create(&exam).Error; err != nil {
		log.Println("[ERROR] Failed to create exam:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create exam"})
	}

	log.Printf("[SUCCESS] Exam created: ID=%d\n", exam.ID)
	return c.Status(201).JSON(fiber.Map{
		"message": "Exam created successfully",
		"data":    exam,
	})
}

// PUT update exam
func (ec *ExamController) UpdateExam(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating exam with ID:", id)

	var exam examModel.ExamModel
	if err := ec.DB.First(&exam, id).Error; err != nil {
		log.Println("[ERROR] Exam not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Exam not found"})
	}

	// Pakai map untuk fleksibel (bisa update sebagian field)
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Optional: bisa tambahkan validasi manual di sini kalau perlu

	if len(updateData) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
	}

	if err := ec.DB.Model(&exam).Updates(updateData).Error; err != nil {
		log.Println("[ERROR] Failed to update exam:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update exam"})
	}

	log.Printf("[SUCCESS] Exam updated: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Exam updated successfully",
		"data":    exam,
	})
}

// DELETE exam
func (ec *ExamController) DeleteExam(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting exam with ID:", id)

	if err := ec.DB.Delete(&examModel.ExamModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete exam:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete exam",
		})
	}

	log.Printf("[SUCCESS] Exam with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Exam with ID %s deleted successfully", id),
	})
}
