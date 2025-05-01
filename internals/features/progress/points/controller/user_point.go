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

func (ctrl *UserPointLogController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var logs []model.UserPointLog
	if err := ctrl.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		log.Println("[ERROR] Gagal mengambil user_point_logs:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data poin user",
		})
	}

	return c.JSON(fiber.Map{
		"data": logs,
	})
}

func (ctrl *UserPointLogController) Create(c *fiber.Ctx) error {
	var input []model.UserPointLog
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
