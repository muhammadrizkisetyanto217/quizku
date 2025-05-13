package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quizku/internals/features/lessons/units/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UnitNewsController struct {
	DB *gorm.DB
}

func NewUnitNewsController(db *gorm.DB) *UnitNewsController {
	return &UnitNewsController{DB: db}
}

// GET all unit news
func (uc *UnitNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.UnitNewsModel
	if err := uc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Unit news list retrieved successfully",
		"data":    news,
	})
}

// GET all news by unit_id
func (uc *UnitNewsController) GetByUnitID(c *fiber.Ctx) error {
	unitID := c.Params("unit_id")

	var news []model.UnitNewsModel
	if err := uc.DB.
		Where("unit_id = ?", unitID).
		Where("deleted_at IS NULL").
		Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if len(news) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "No news found for this unit_id",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News for the selected unit retrieved successfully",
		"data":    news,
	})
}

// GET by ID
func (uc *UnitNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unit news found successfully",
		"data":    news,
	})
}

// CREATE
func (uc *UnitNewsController) Create(c *fiber.Ctx) error {
	var news model.UnitNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := uc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Unit news created successfully",
		"data":    news,
	})
}

// UPDATE
func (uc *UnitNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := uc.DB.Save(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.JSON(fiber.Map{
		"message": "Unit news updated successfully",
		"data":    news,
	})
}

// DELETE
func (uc *UnitNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	if err := uc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Unit news with ID %v deleted successfully", news.ID),
	})
}

// Optional: helper update JSON field (kalau units punya kolom update_news)
func updateUnitNewsJSON(db *gorm.DB, unitID int) {
	var newsList []model.UnitNewsModel
	if err := db.Where("unit_id = ?", unitID).Order("created_at desc").Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch unit news for update:", err)
		return
	}

	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal unit news:", err)
		return
	}

	// Pastikan kolom update_news sudah ada di tabel units
	db.Table("units").
		Where("id = ?", unitID).
		Update("update_news", datatypes.JSON(newsData))
}