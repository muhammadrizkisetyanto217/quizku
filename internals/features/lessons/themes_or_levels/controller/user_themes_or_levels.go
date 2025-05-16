package controller

import (
	"log"

	themesOrLevelsModel "quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserThemesController struct {
	DB *gorm.DB
}

func NewUserThemesController(db *gorm.DB) *UserThemesController {
	return &UserThemesController{DB: db}
}

// 🟢 GET /api/user-themes/:user_id
// Mengambil semua data progres themes_or_levels yang dimiliki oleh user tertentu berdasarkan user_id.
// Digunakan untuk menampilkan rekap progres per tema seperti: unit lengkap, grade, dan status kelulusan.
func (ctrl *UserThemesController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	// Validasi UUID dari parameter
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []themesOrLevelsModel.UserThemesOrLevelsModel

	// Ambil semua data user_themes_or_levels milik user tersebut
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_themes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	// Kirim response dengan data progres user per theme
	return c.JSON(fiber.Map{
		"data": data,
	})
}
