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

func (ctrl *UserSectionQuizzesController) GetUserSectionQuizzesByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []model.UserSectionQuizzesModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil user_section_quizzes:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data user_section_quizzes",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
