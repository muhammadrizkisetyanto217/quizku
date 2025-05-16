package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"quizku/internals/features/progress/points/model" // sesuaikan path
)

type UserPointLogController struct {
	DB *gorm.DB
}

func NewUserPointLogController(db *gorm.DB) *UserPointLogController {
	return &UserPointLogController{DB: db}
}

// 🟢 GET /api/user-point-logs/:user_id
// Mengambil seluruh riwayat poin milik user berdasarkan user_id.
// Digunakan untuk menampilkan log aktivitas user seperti kuis, evaluasi, reading, dsb.
func (ctrl *UserPointLogController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	// Validasi UUID
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var logs []model.UserPointLog

	// Ambil semua log user yang sesuai user_id, diurutkan dari terbaru
	if err := ctrl.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {

		log.Println("[ERROR] Gagal mengambil user_point_logs:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data poin user",
		})
	}

	// Kirim data log ke client
	return c.JSON(fiber.Map{
		"data": logs,
	})
}

// 🟡 POST /api/user-point-logs
// Menambahkan banyak log poin dalam sekali kirim (batch).
// Cocok untuk digunakan oleh service poin dari quiz, reading, exam, dsb.
func (ctrl *UserPointLogController) Create(c *fiber.Ctx) error {
	var input []model.UserPointLog

	// Validasi format body
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Body parser gagal:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format input tidak valid",
		})
	}

	if len(input) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data logs kosong",
		})
	}

	// Simpan ke database
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Gagal menyimpan logs:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menyimpan data poin",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil menambahkan log poin",
		"count":   len(input),
	})
}
