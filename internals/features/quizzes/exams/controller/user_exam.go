package controller

import (
	"log"
	"net/http"
	"time"

	"quizku/internals/features/quizzes/exams/model"
	"quizku/internals/features/quizzes/exams/service"

	"github.com/go-playground/validator/v10"
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

// Create user_exam
func (c *UserExamController) Create(ctx *fiber.Ctx) error {
	// ✅ Ambil user_id dari JWT
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

	// ✅ Struct input body tanpa user_id
	type InputBody struct {
		ExamID          uint    `json:"exam_id" validate:"required"`
		UnitID          uint    `json:"unit_id" validate:"required"`
		PercentageGrade float64 `json:"percentage_grade" validate:"required"`
		TimeDuration    int     `json:"time_duration"`
		Point           int     `json:"point"`
	}

	var body InputBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// ✅ Validasi
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing or invalid fields",
			"error":   err.Error(),
		})
	}

	// ✅ Cek jika sudah ada entry sebelumnya
	var existing model.UserExamModel
	err = c.DB.Where("user_id = ? AND exam_id = ?", userUUID, body.ExamID).
		First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("[ERROR] Gagal cek user_exam existing:", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memproses data",
			"error":   err.Error(),
		})
	}

	if err == nil {
		// ✅ Sudah ada → update attempt dan nilai jika lebih tinggi
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

		_ = service.AddPointFromExam(c.DB, existing.UserID, existing.ExamID, existing.Attempt)

		// ✅ Tambahkan pencatatan aktivitas harian
		_ = activityService.UpdateOrInsertDailyActivity(c.DB, existing.UserID)

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"message": "User exam record updated successfully",
			"data":    existing,
		})
	}

	// ✅ Belum ada → buat baru
	newExam := model.UserExamModel{
		UserID:          userUUID,
		ExamID:          body.ExamID,
		UnitID:          body.UnitID,
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

	_ = service.AddPointFromExam(c.DB, newExam.UserID, newExam.ExamID, newExam.Attempt)

	// ✅ Tambahkan pencatatan aktivitas harian
	_ = activityService.UpdateOrInsertDailyActivity(c.DB, newExam.UserID)

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User exam record created successfully",
		"data":    newExam,
	})
}

// Delete user_exam by ID
func (c *UserExamController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var exam model.UserExamModel
	if err := c.DB.First(&exam, id).Error; err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}

	if err := c.DB.Delete(&exam).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user exam",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User exam deleted successfully",
	})
}

// Get all user_exams
func (c *UserExamController) GetAll(ctx *fiber.Ctx) error {
	var data []model.UserExamModel
	if err := c.DB.Find(&data).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve data",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

// Get user_exam by ID
func (c *UserExamController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var data model.UserExamModel
	if err := c.DB.First(&data, id).Error; err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

// Get user_exams by user_id (UUID)
func (ctrl *UserExamController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []model.UserExamModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_exam:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
