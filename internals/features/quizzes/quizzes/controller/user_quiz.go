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
//
// Langkah-langkah:
//   - Ambil user_id dari token JWT
//   - Validasi input: quiz_id dan percentage_grade wajib
//   - Cek apakah user sudah pernah mengerjakan quiz
//     🔸 Jika belum → buat entri baru dengan attempt = 1
//     🔸 Jika sudah → update record, tambahkan attempt, simpan grade terbaik
//   - Ambil relasi dari quiz → section → unit
//   - Jalankan logika progres:
//     🔹 Update progress section jika quiz lengkap
//     🔹 Update progress unit jika semua section lengkap
//     🔹 Tambahkan poin
func (uc *UserQuizController) CreateOrUpdateUserQuiz(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user quiz progress")

	// ✅ Ambil user_id dari JWT
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// ✅ Struct input dari body
	type InputBody struct {
		QuizID          uint `json:"quiz_id" validate:"required"`
		PercentageGrade int  `json:"percentage_grade" validate:"required"`
		TimeDuration    int  `json:"time_duration"` // opsional
		Point           int  `json:"point"`         // opsional
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse input:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// ✅ Validasi field wajib
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// ✅ Cek apakah user sudah punya data kuis
	var existing model.UserQuizzesModel
	err = uc.DB.Where("user_id = ? AND quiz_id = ?", userUUID, body.QuizID).First(&existing).Error

	var attempt int
	var finalGrade int

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 🔸 Buat data baru
		attempt = 1
		finalGrade = body.PercentageGrade
		newRecord := model.UserQuizzesModel{
			UserID:          userUUID,
			QuizID:          body.QuizID,
			Attempt:         attempt,
			PercentageGrade: finalGrade,
			TimeDuration:    body.TimeDuration,
			Point:           body.Point,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := uc.DB.Create(&newRecord).Error; err != nil {
			log.Println("[ERROR] Failed to create user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user quiz"})
		}
		existing = newRecord
		log.Printf("[SUCCESS] Created user_quiz for user_id=%s quiz_id=%d\n", userUUID, body.QuizID)

	} else if err == nil {
		// 🔸 Update existing
		attempt = existing.Attempt + 1
		finalGrade = max(existing.PercentageGrade, body.PercentageGrade)

		existing.Attempt = attempt
		existing.PercentageGrade = finalGrade
		existing.TimeDuration = body.TimeDuration
		existing.Point = body.Point
		existing.UpdatedAt = time.Now()

		if err := uc.DB.Save(&existing).Error; err != nil {
			log.Println("[ERROR] Failed to update user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user quiz"})
		}
		log.Printf("[SUCCESS] Updated user_quiz (attempt %d, grade %d) for user_id=%s quiz_id=%d\n",
			attempt, finalGrade, userUUID, body.QuizID)

	} else {
		log.Println("[ERROR] Failed to fetch user quiz:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user quiz"})
	}

	// ✅ Ambil struktur relasi quiz → section → unit
	var quiz model.QuizModel
	if err := uc.DB.First(&quiz, body.QuizID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Quiz not found"})
	}
	var section model.SectionQuizzesModel
	if err := uc.DB.First(&section, quiz.SectionQuizID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Section not found"})
	}
	var unit unitModel.UnitModel
	if err := uc.DB.First(&unit, section.UnitID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch unit"})
	}

	// ✅ Update progress section → unit → poin
	_ = services.UpdateUserSectionIfQuizCompleted(uc.DB, userUUID, section.ID, body.QuizID, attempt, finalGrade)
	_ = services.UpdateUserUnitIfSectionCompleted(uc.DB, userUUID, unit.ID, section.ID)
	if err := services.AddPointFromQuiz(uc.DB, userUUID, body.QuizID, attempt); err != nil {
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

	// ✅ Validasi format UUID
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	// ✅ Ambil seluruh data user_quizzes milik user
	var userQuizzes []model.UserQuizzesModel
	if err := uc.DB.Where("user_id = ?", userID).Find(&userQuizzes).Error; err != nil {
		log.Println("[ERROR] Gagal mengambil user_quizzes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data quiz user",
		})
	}

	// ✅ Kembalikan respons
	return c.JSON(fiber.Map{
		"data": userQuizzes,
	})
}
