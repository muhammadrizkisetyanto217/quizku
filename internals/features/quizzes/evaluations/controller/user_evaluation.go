package controller

import (
	"log"
	userEvaluationModel "quizku/internals/features/quizzes/evaluations/model"
	"quizku/internals/features/quizzes/evaluations/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	activityService "quizku/internals/features/progress/daily_activities/service"
)

type UserEvaluationController struct {
	DB *gorm.DB
}

func NewUserEvaluationController(db *gorm.DB) *UserEvaluationController {
	return &UserEvaluationController{DB: db}
}

// POST /api/user_evaluations3
func (ctrl *UserEvaluationController) Create(c *fiber.Ctx) error {
	// ✅ Ambil user_id dari JWT
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// ✅ Parse body tanpa user_id
	type InputBody struct {
		EvaluationID    uint `json:"evaluation_id" validate:"required"`
		UnitID          uint `json:"unit_id" validate:"required"`
		PercentageGrade int  `json:"percentage_grade" validate:"required"`
		TimeDuration    int  `json:"time_duration"` // opsional
		Point           int  `json:"point"`         // opsional
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ Validasi
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// ✅ Hitung attempt ke-n
	var latestAttempt int
	err = ctrl.DB.Table("user_evaluations").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND evaluation_id = ?", userUUID, body.EvaluationID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Siapkan data akhir
	input := userEvaluationModel.UserEvaluationModel{
		UserID:          userUUID,
		EvaluationID:    body.EvaluationID,
		UnitID:          body.UnitID,
		Attempt:         latestAttempt + 1,
		PercentageGrade: body.PercentageGrade,
		TimeDuration:    body.TimeDuration,
		Point:           body.Point,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// ✅ Simpan ke DB
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user evaluation:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user evaluation"})
	}

	// ✅ Update progres & poin
	if err := service.UpdateUserUnitFromEvaluation(ctrl.DB, input.UserID, input.UnitID, input.PercentageGrade); err != nil {
		log.Println("[ERROR] Gagal update user_unit:", err)
	}
	if err := service.AddPointFromEvaluation(ctrl.DB, input.UserID, input.EvaluationID, input.Attempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin:", err)
	}

	// ✅ Tambahkan aktivitas harian
	if err := activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserID); err != nil {
		log.Println("[ERROR] Gagal mencatat aktivitas harian:", err)
	}

	log.Printf("[SUCCESS] UserEvaluation created: user_id=%s, evaluation_id=%d, attempt=%d\n",
		input.UserID.String(), input.EvaluationID, input.Attempt)

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
