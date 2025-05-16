package controller

import (
	"log"
	"net/http"
	"time"

	"quizku/internals/features/quizzes/exams/model"
	examModel "quizku/internals/features/quizzes/exams/model"
	"quizku/internals/features/quizzes/exams/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	activityService "quizku/internals/features/progress/daily_activities/service"
)

type UserExamController struct {
	DB *gorm.DB
}

func NewUserExamController(db *gorm.DB) *UserExamController {
	return &UserExamController{DB: db}
}

// 🟡 POST /api/user-exams
// Menyimpan atau memperbarui hasil ujian (exam) yang dikerjakan oleh user.
// Jika user sudah pernah mengerjakan exam → data akan di-*update* dan attempt bertambah.
// Jika belum pernah → akan membuat record baru.
//
// Fungsi ini otomatis:
// ✅ Menambahkan attempt ke-n,
// ✅ Menyimpan poin dan nilai tertinggi,
// ✅ Menyimpan durasi pengerjaan,
// ✅ Memperbarui aktivitas harian,
// ✅ Menambahkan poin ke user_point_log.
func (c *UserExamController) Create(ctx *fiber.Ctx) error {
	// 🔐 Ambil user_id dari token JWT
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid user ID format",
		})
	}

	// 📥 Struktur body request
	type InputBody struct {
		ExamID          uint    `json:"exam_id"`          // ID exam
		PercentageGrade float64 `json:"percentage_grade"` // Nilai (0-100)
		TimeDuration    int     `json:"time_duration"`    // Lama waktu pengerjaan dalam detik
		Point           int     `json:"point"`            // Poin yang didapat
	}
	var body InputBody

	// ✅ Parse dan validasi body
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if body.ExamID == 0 || body.PercentageGrade == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "exam_id and percentage_grade are required",
		})
	}

	// 🔎 Ambil unit_id dari exam yang terkait
	var exam examModel.ExamModel
	if err := c.DB.Select("id, unit_id").First(&exam, body.ExamID).Error; err != nil {
		log.Println("[ERROR] Exam not found:", err)
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Exam not found",
		})
	}

	// 🔄 Cek apakah user sudah pernah mengerjakan exam ini
	var existing model.UserExamModel
	err = c.DB.Where("user_id = ? AND exam_id = ?", userUUID, body.ExamID).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("[ERROR] Gagal cek user_exam existing:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memproses data",
			"error":   err.Error(),
		})
	}

	// 🔁 Jika sudah ada → update nilai & tambah attempt
	if err == nil {
		existing.Attempt += 1
		if body.PercentageGrade > float64(existing.PercentageGrade) {
			existing.PercentageGrade = int(body.PercentageGrade)
		}
		existing.TimeDuration = body.TimeDuration
		existing.Point = body.Point

		if err := c.DB.Save(&existing).Error; err != nil {
			log.Println("[ERROR] Gagal update user_exam:", err)
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal memperbarui data",
				"error":   err.Error(),
			})
		}

		// ✅ Tambah poin & aktivitas harian
		_ = service.AddPointFromExam(c.DB, existing.UserID, existing.ExamID, existing.Attempt)
		_ = activityService.UpdateOrInsertDailyActivity(c.DB, existing.UserID)

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"message": "User exam record updated successfully",
			"data":    existing,
		})
	}

	// ✅ Jika belum pernah → buat baru
	newExam := model.UserExamModel{
		UserID:          userUUID,
		ExamID:          body.ExamID,
		UnitID:          exam.UnitID,
		Attempt:         1,
		PercentageGrade: int(body.PercentageGrade),
		TimeDuration:    body.TimeDuration,
		Point:           body.Point,
		CreatedAt:       time.Now(),
	}

	if err := c.DB.Create(&newExam).Error; err != nil {
		log.Println("[ERROR] Gagal create user_exam:", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user exam record",
			"error":   err.Error(),
		})
	}

	// ✅ Tambah poin & aktivitas harian
	_ = service.AddPointFromExam(c.DB, newExam.UserID, newExam.ExamID, newExam.Attempt)
	_ = activityService.UpdateOrInsertDailyActivity(c.DB, newExam.UserID)

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User exam record created successfully",
		"data":    newExam,
	})
}
// 🔴 DELETE /api/user-exams/:id
// Menghapus 1 data `user_exam` berdasarkan ID (bukan UUID user).
// Cocok digunakan saat admin ingin membatalkan hasil exam user tertentu.
func (c *UserExamController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var exam model.UserExamModel
	if err := c.DB.First(&exam, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}

	if err := c.DB.Delete(&exam).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user exam",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User exam deleted successfully",
	})
}

// 🟢 GET /api/user-exams
// Mengambil seluruh data user_exam tanpa filter.
// Cocok untuk keperluan debug, export, atau admin monitoring.
func (c *UserExamController) GetAll(ctx *fiber.Ctx) error {
	var data []model.UserExamModel
	if err := c.DB.Find(&data).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve data",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

// 🟢 GET /api/user-exams/:id
// Mengambil satu data user_exam berdasarkan ID record.
// Biasanya digunakan untuk halaman detail hasil ujian.
func (c *UserExamController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var data model.UserExamModel
	if err := c.DB.First(&data, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

// 🟢 GET /api/user-exams/user/:user_id
// Mengambil seluruh hasil exam milik satu user (berdasarkan UUID user).
// Cocok digunakan untuk halaman riwayat ujian user atau profil.
func (ctrl *UserExamController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []model.UserExamModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_exam:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
