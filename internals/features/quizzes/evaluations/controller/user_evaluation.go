package controller

import (
	"log"
	userEvaluationModel "quizku/internals/features/quizzes/evaluations/model"
	"quizku/internals/features/quizzes/evaluations/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserEvaluationController struct {
	DB *gorm.DB
}

func NewUserEvaluationController(db *gorm.DB) *UserEvaluationController {
	return &UserEvaluationController{DB: db}
}

// POST /api/user_evaluations3
func (ctrl *UserEvaluationController) Create(c *fiber.Ctx) error {
	var input userEvaluationModel.UserEvaluationModel
	body := c.Body()
	log.Println("[DEBUG] Raw request body:", string(body))

	// ✅ Parse body
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// ✅ Hitung attempt ke-n
	var latestAttempt int
	err := ctrl.DB.Table("user_evaluations").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND evaluation_id = ?", input.UserID, input.EvaluationID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	input.Attempt = latestAttempt + 1

	// ✅ Simpan ke DB
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user evaluation:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user evaluation"})
	}

	// ✅ Update ke user_unit hanya jika grade lebih tinggi
	if err := service.UpdateUserUnitFromEvaluation(ctrl.DB, input.UserID, input.UnitID, input.PercentageGrade); err != nil {
		log.Println("[ERROR] Failed to update user unit from evaluation:", err)
	}

	// Tambahkan log point dari evaluation
	if err := service.AddPointFromEvaluation(ctrl.DB, input.UserID, input.EvaluationID, input.Attempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin dari evaluation:", err)
	}

	log.Println("[SUCCESS] User evaluation created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User evaluation created successfully",
		"data":    input,
	})
}

// GET /api/user_evaluations/:user_id
func (ctrl *UserEvaluationController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var evaluations []userEvaluationModel.UserEvaluationModel

	if err := ctrl.DB.Where("user_id = ?", userID).Find(&evaluations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get evaluations"})
	}

	return c.JSON(evaluations)
}
