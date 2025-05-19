package controller

import (
	"errors"
	"log"
	"time"

	"quizku/internals/features/quizzes/quizzes/model"
	"quizku/internals/features/quizzes/quizzes/services"

	unitModel "quizku/internals/features/lessons/units/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserQuizController struct {
	DB *gorm.DB
}

func NewUserQuizController(db *gorm.DB) *UserQuizController {
	return &UserQuizController{DB: db}
}

// ✅ POST /api/user-quizzes
// Membuat atau memperbarui progres pengerjaan kuis oleh user, sekaligus mengatur progres section dan unit.
func (uc *UserQuizController) CreateOrUpdateUserQuiz(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user quiz progress")

	// 🔐 Ambil user_id dari JWT
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// 📥 Parsing input (pakai nama field semantik)
	type InputBody struct {
		UserQuizQuizID          uint `json:"user_quiz_quiz_id" validate:"required"`
		UserQuizPercentageGrade int  `json:"user_quiz_percentage_grade" validate:"required"`
		UserQuizTimeDuration    int  `json:"user_quiz_time_duration"`
		UserQuizPoint           int  `json:"user_quiz_point"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse input:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// 🧠 Cek apakah user sudah pernah mengerjakan kuis tersebut
	var existing model.UserQuizzesModel
	err = uc.DB.Where("user_quiz_user_id = ? AND user_quiz_quiz_id = ?", userUUID, body.UserQuizQuizID).
		First(&existing).Error

	var userQuizAttempt int
	var userQuizBestScore int

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 🔹 Belum pernah → buat baru
		userQuizAttempt = 1
		userQuizBestScore = body.UserQuizPercentageGrade

		newRecord := model.UserQuizzesModel{
			UserQuizUserID:          userUUID,
			UserQuizQuizID:          body.UserQuizQuizID,
			UserQuizAttempt:         userQuizAttempt,
			UserQuizPercentageGrade: userQuizBestScore,
			UserQuizTimeDuration:    body.UserQuizTimeDuration,
			UserQuizPoint:           body.UserQuizPoint,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}
		if err := uc.DB.Create(&newRecord).Error; err != nil {
			log.Println("[ERROR] Failed to create user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user quiz"})
		}
		existing = newRecord
		log.Printf("[SUCCESS] Created user_quiz for user_id=%s quiz_id=%d\n", userUUID, body.UserQuizQuizID)

	} else if err == nil {
		// 🔹 Sudah pernah → update
		userQuizAttempt = existing.UserQuizAttempt + 1
		userQuizBestScore = max(existing.UserQuizPercentageGrade, body.UserQuizPercentageGrade)

		existing.UserQuizAttempt = userQuizAttempt
		existing.UserQuizPercentageGrade = userQuizBestScore
		existing.UserQuizTimeDuration = body.UserQuizTimeDuration
		existing.UserQuizPoint = body.UserQuizPoint
		existing.UpdatedAt = time.Now()

		if err := uc.DB.Save(&existing).Error; err != nil {
			log.Println("[ERROR] Failed to update user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user quiz"})
		}
		log.Printf("[SUCCESS] Updated user_quiz (attempt %d, grade %d) for user_id=%s quiz_id=%d\n",
			userQuizAttempt, userQuizBestScore, userUUID, body.UserQuizQuizID)

	} else {
		log.Println("[ERROR] Failed to fetch user quiz:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user quiz"})
	}

	// 🔍 Ambil relasi quiz → section → unit
	var quiz model.QuizModel
	if err := uc.DB.First(&quiz, "quiz_id = ?", body.UserQuizQuizID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Quiz not found"})
	}

	var section model.SectionQuizzesModel
	if err := uc.DB.First(&section, "section_quizzes_id = ?", quiz.QuizSectionQuizzesID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Section not found"})
	}

	var unit unitModel.UnitModel
	if err := uc.DB.First(&unit, "unit_id = ?", section.SectionQuizzesUnitID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch unit"})
	}

	// 🔁 Update progres section
	_ = services.UpdateUserSectionIfQuizCompleted(
		uc.DB,
		userUUID,                 // user_id dari JWT
		section.SectionQuizzesID, // section tempat quiz berada
		body.UserQuizQuizID,      // quiz yang dikerjakan
		userQuizAttempt,          // attempt user
		userQuizBestScore,        // grade tertinggi user
	)

	// 🔁 Update progres unit jika semua section selesai
	_ = services.UpdateUserUnitIfSectionCompleted(
		uc.DB,
		userUUID,                 // userID: UUID user yang login
		unit.UnitID,              // unitID: ID unit terkait
		section.SectionQuizzesID, // completedSectionID: ID section yang baru selesai
	)

	// ➕ Tambah poin dari quiz
	if err := services.AddPointFromQuiz(uc.DB, userUUID, body.UserQuizQuizID, userQuizAttempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin dari quiz:", err)
	}

	return c.JSON(fiber.Map{
		"message": "User quiz progress saved and progress updated",
		"data":    existing,
	})
}

// Helper untuk nilai maksimum (bisa disimpan global)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ✅ GET /api/user-quizzes/:user_id
// Mengambil seluruh riwayat pengerjaan kuis oleh user berdasarkan user_id.
//
// Langkah-langkah:
// - Ambil dan validasi parameter user_id (UUID format)
// - Query tabel user_quizzes berdasarkan user_id
// - Kembalikan data list pengerjaan kuis (termasuk attempt, grade, dan timestamp)
func (uc *UserQuizController) GetUserQuizzesByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var userQuizzes []model.UserQuizzesModel
	if err := uc.DB.
		Where("user_quiz_user_id = ?", userID).
		Order("updated_at DESC").
		Find(&userQuizzes).Error; err != nil {

		log.Println("[ERROR] Gagal mengambil user_quizzes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data quiz user",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(userQuizzes),
		"data":  userQuizzes,
	})
}
