package controller

import (
	"log"
	evaluationModel "quizku/internals/features/quizzes/evaluations/model"
	userEvaluationModel "quizku/internals/features/quizzes/evaluations/model"
	"quizku/internals/features/quizzes/evaluations/service"
	"time"

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

// 🟡 POST /api/user_evaluations3
// Menyimpan hasil pengerjaan evaluasi oleh user (attempt).
// Fungsi ini otomatis:
// - Mengisi attempt ke-n (berdasarkan data sebelumnya),
// - Mengupdate progress di user_unit,
// - Menambahkan poin ke user_point_log,
// - Mencatat aktivitas harian user.
func (ctrl *UserEvaluationController) Create(c *fiber.Ctx) error {
	// 🔐 Ambil user_id dari token (middleware auth)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// 📦 Struktur input JSON sesuai model
	type InputBody struct {
		UserEvaluationEvaluationID    uint `json:"user_evaluation_evaluation_id"`
		UserEvaluationPercentageGrade int  `json:"user_evaluation_percentage_grade"`
		UserEvaluationTimeDuration    int  `json:"user_evaluation_time_duration"`
		UserEvaluationPoint           int  `json:"user_evaluation_point"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if body.UserEvaluationEvaluationID == 0 || body.UserEvaluationPercentageGrade == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_evaluation_evaluation_id and user_evaluation_percentage_grade are required",
		})
	}

	// 🔎 Ambil evaluasi untuk ambil unit_id
	var evaluation evaluationModel.EvaluationModel
	if err := ctrl.DB.
		Select("evaluation_id, evaluation_unit_id").
		First(&evaluation, "evaluation_id = ?", body.UserEvaluationEvaluationID).Error; err != nil {
		log.Println("[ERROR] Evaluation not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Evaluation not found"})
	}

	// 🔁 Ambil attempt terakhir dari user untuk evaluasi ini
	var latestAttempt int
	err = ctrl.DB.Table("user_evaluations").
		Select("COALESCE(MAX(user_evaluation_attempt), 0)").
		Where("user_evaluation_user_id = ? AND user_evaluation_evaluation_id = ?", userUUID, body.UserEvaluationEvaluationID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// 🧾 Buat entri baru
	input := userEvaluationModel.UserEvaluationModel{
		UserEvaluationUserID:          userUUID,
		UserEvaluationEvaluationID:    body.UserEvaluationEvaluationID,
		UserEvaluationUnitID:          evaluation.EvaluationUnitID,
		UserEvaluationAttempt:         latestAttempt + 1,
		UserEvaluationPercentageGrade: body.UserEvaluationPercentageGrade,
		UserEvaluationTimeDuration:    body.UserEvaluationTimeDuration,
		UserEvaluationPoint:           body.UserEvaluationPoint,
		CreatedAt:                     time.Now(),
		UpdatedAt:                     time.Now(),
	}

	// 💾 Simpan ke database
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user evaluation:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user evaluation"})
	}

	// ⛏️ Update progres, poin, dan aktivitas harian
	_ = service.UpdateUserUnitFromEvaluation(ctrl.DB, input.UserEvaluationUserID, input.UserEvaluationUnitID, input.UserEvaluationPercentageGrade)
	_ = service.AddPointFromEvaluation(ctrl.DB, input.UserEvaluationUserID, input.UserEvaluationEvaluationID, input.UserEvaluationAttempt)
	_ = activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserEvaluationUserID)

	log.Printf("[SUCCESS] UserEvaluation created: user_id=%s, evaluation_id=%d, attempt=%d\n",
		input.UserEvaluationUserID.String(), input.UserEvaluationEvaluationID, input.UserEvaluationAttempt)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User evaluation created successfully",
		"data":    input,
	})
}

// 🟢 GET /api/user_evaluations/:user_id
// Mengambil seluruh data evaluasi yang sudah pernah dikerjakan oleh user berdasarkan user_id.
func (ctrl *UserEvaluationController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var evaluations []userEvaluationModel.UserEvaluationModel

	log.Printf("[INFO] Fetching user evaluations for user_id: %s", userID)

	// 🔍 Ambil semua evaluasi berdasarkan user_evaluation_user_id
	if err := ctrl.DB.
		Where("user_evaluation_user_id = ?", userID).
		Find(&evaluations).Error; err != nil {
		log.Println("[ERROR] Failed to get user evaluations:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get evaluations",
		})
	}

	log.Printf("[SUCCESS] Found %d user evaluations for user_id: %s", len(evaluations), userID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User evaluations fetched successfully",
		"total":   len(evaluations),
		"data":    evaluations,
	})
}
