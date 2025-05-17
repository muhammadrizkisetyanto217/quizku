package controller

import (
	"log"
	"net/http"

	"quizku/internals/features/quizzes/quizzes/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSectionQuizzesController struct {
	DB *gorm.DB
}

func NewUserSectionQuizzesController(db *gorm.DB) *UserSectionQuizzesController {
	return &UserSectionQuizzesController{
		DB: db,
	}
}

// ✅ GET /api/user-section-quizzes/:user_id
// Mengambil seluruh progres section quiz yang telah dikerjakan oleh user berdasarkan user_id.
//
// Langkah-langkah:
// - Ambil user_id dari path parameter
// - Validasi bahwa user_id merupakan UUID yang valid
// - Query ke tabel user_section_quizzes untuk mendapatkan semua progres milik user tersebut
// - Kembalikan dalam bentuk array JSON
func (ctrl *UserSectionQuizzesController) GetUserSectionQuizzesByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")

	// ✅ Validasi UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	// ✅ Ambil data user_section_quizzes dari database
	var data []model.UserSectionQuizzesModel
	if err := ctrl.DB.
		Where("user_section_quizzes_user_id = ?", userID).
		Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil user_section_quizzes:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data user_section_quizzes",
		})
	}

	// ✅ Kembalikan hasil sebagai respons
	return c.JSON(fiber.Map{
		"data": data,
	})
}
