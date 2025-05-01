package controller

import (
	"log"
	"net/http"

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

// GET /api/user-themes/:user_id
func (ctrl *UserThemesController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_themes:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

