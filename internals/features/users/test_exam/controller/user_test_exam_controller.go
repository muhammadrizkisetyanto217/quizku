package controller

import (
	"quizku/internals/features/users/test_exam/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTestExamController struct {
	DB *gorm.DB
}

func NewUserTestExamController(db *gorm.DB) *UserTestExamController {
	return &UserTestExamController{DB: db}
}

// ✅ Ambil semua data user_test_exam
func (ctrl *UserTestExamController) GetAll(c *fiber.Ctx) error {
	var results []model.UserTestExam
	if err := ctrl.DB.Order("user_test_exam_id DESC").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data user_test_exam"})
	}
	return c.JSON(results)
}

// ✅ Ambil satu user_test_exam berdasarkan ID
func (ctrl *UserTestExamController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var data model.UserTestExam
	if err := ctrl.DB.First(&data, "user_test_exam_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User test exam tidak ditemukan"})
	}
	return c.JSON(data)
}

// ✅ Buat entri user_test_exam baru
func (ctrl *UserTestExamController) Create(c *fiber.Ctx) error {
	var payload model.UserTestExam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	// Validasi user_id & test_exam_id tidak boleh kosong
	if payload.UserTestExamUserID == uuid.Nil || payload.UserTestExamTestExamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID dan Test Exam ID wajib diisi",
		})
	}

	payload.CreatedAt = time.Now()

	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User test exam berhasil disimpan",
		"data":    payload,
	})
}

// ✅ Ambil semua hasil test exam berdasarkan user_id
func (ctrl *UserTestExamController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var results []model.UserTestExam
	if err := ctrl.DB.
		Where("user_test_exam_user_id = ?", userID).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data berdasarkan user_id",
		})
	}
	return c.JSON(results)
}

// ✅ Ambil semua hasil peserta untuk test_exam tertentu
func (ctrl *UserTestExamController) GetByTestExamID(c *fiber.Ctx) error {
	examID := c.Params("test_exam_id")
	var results []model.UserTestExam
	if err := ctrl.DB.
		Where("user_test_exam_test_exam_id = ?", examID).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data berdasarkan test_exam_id",
		})
	}
	return c.JSON(results)
}
