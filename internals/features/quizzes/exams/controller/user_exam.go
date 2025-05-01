package controller

import (
	"log"
	"net/http"

	"quizku/internals/features/quizzes/exams/model"
	"quizku/internals/features/quizzes/exams/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserExamController struct {
	DB *gorm.DB
}

func NewUserExamController(db *gorm.DB) *UserExamController {
	return &UserExamController{DB: db}
}

// Create user_exam
func (c *UserExamController) Create(ctx *fiber.Ctx) error {
	var payload model.UserExamModel

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi user_id
	if payload.UserID == uuid.Nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "UserID is required and must be a valid UUID",
		})
	}

	// Cek apakah sudah ada user_exam untuk kombinasi user_id + exam_id
	var existing model.UserExamModel
	err := c.DB.Where("user_id = ? AND exam_id = ?", payload.UserID, payload.ExamID).
		First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("[ERROR] Gagal cek user_exam existing:", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memproses data",
			"error":   err.Error(),
		})
	}

	if err == nil {
		// Sudah ada → update (attempt++, nilai tertinggi)
		existing.Attempt += 1
		if payload.PercentageGrade > existing.PercentageGrade {
			existing.PercentageGrade = payload.PercentageGrade
		}

		if err := c.DB.Save(&existing).Error; err != nil {
			log.Println("[ERROR] Gagal update user_exam:", err)
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal memperbarui data",
				"error":   err.Error(),
			})
		}

		// Tambahkan log point
		_ = service.AddPointFromExam(c.DB, existing.UserID, existing.ExamID, existing.Attempt)

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"message": "User exam record updated successfully",
			"data":    existing,
		})
	}

	// Belum ada → buat baru
	payload.Attempt = 1
	if err := c.DB.Create(&payload).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user exam record",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User exam record created successfully",
		"data":    payload,
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