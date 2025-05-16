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

	// 📦 Struktur input body
	type InputBody struct {
		EvaluationID    uint `json:"evaluation_id"`    // ID evaluasi yang dikerjakan
		PercentageGrade int  `json:"percentage_grade"` // Nilai persentase
		TimeDuration    int  `json:"time_duration"`    // Lama waktu pengerjaan (detik)
		Point           int  `json:"point"`            // Poin yang didapatkan dari evaluasi ini
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if body.EvaluationID == 0 || body.PercentageGrade == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "evaluation_id and percentage_grade are required",
		})
	}

	// 🔎 Ambil evaluasi dan unit_id-nya
	var evaluation evaluationModel.EvaluationModel
	if err := ctrl.DB.Select("id, unit_id").First(&evaluation, body.EvaluationID).Error; err != nil {
		log.Println("[ERROR] Evaluation not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Evaluation not found"})
	}

	// 🔁 Cek attempt terakhir user untuk evaluasi ini
	var latestAttempt int
	err = ctrl.DB.Table("user_evaluations").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND evaluation_id = ?", userUUID, body.EvaluationID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// 📤 Siapkan data untuk disimpan
	input := userEvaluationModel.UserEvaluationModel{
		UserID:          userUUID,
		EvaluationID:    body.EvaluationID,
		UnitID:          evaluation.UnitID,
		Attempt:         latestAttempt + 1,
		PercentageGrade: body.PercentageGrade,
		TimeDuration:    body.TimeDuration,
		Point:           body.Point,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 💾 Simpan ke DB
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user evaluation:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user evaluation"})
	}

	// ⛏️ Update progress & poin
	_ = service.UpdateUserUnitFromEvaluation(ctrl.DB, input.UserID, input.UnitID, input.PercentageGrade)
	_ = service.AddPointFromEvaluation(ctrl.DB, input.UserID, input.EvaluationID, input.Attempt)
	_ = activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserID)

	log.Printf("[SUCCESS] UserEvaluation created: user_id=%s, evaluation_id=%d, attempt=%d\n",
		input.UserID.String(), input.EvaluationID, input.Attempt)

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

	// Ambil seluruh log evaluasi milik user
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&evaluations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get evaluations"})
	}

	return c.JSON(evaluations)
}
