package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ThemesOrLevelsNewsController struct {
	DB *gorm.DB
}

func NewThemesOrLevelsNewsController(db *gorm.DB) *ThemesOrLevelsNewsController {
	return &ThemesOrLevelsNewsController{DB: db}
}

// GET all
func (tc *ThemesOrLevelsNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.ThemesOrLevelsNewsModel
	if err := tc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Themes/Levels news list retrieved successfully",
		"data":    news,
	})
}

// GET by ThemesOrLevelsID
func (tc *ThemesOrLevelsNewsController) GetByThemesOrLevelsID(c *fiber.Ctx) error {
	id := c.Params("themes_or_levels_id") // themes_or_levels_id dari URL
	var news []model.ThemesOrLevelsNewsModel

	if err := tc.DB.
		Where("themes_or_levels_id = ?", id).
		Where("deleted_at IS NULL").
		Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if len(news) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "No news found for this themes_or_levels_id",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News for the selected themes/levels retrieved successfully",
		"data":    news,
	})
}

// GET by ID
func (tc *ThemesOrLevelsNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Themes/Levels news found successfully",
		"data":    news,
	})
}

// CREATE
func (tc *ThemesOrLevelsNewsController) Create(c *fiber.Ctx) error {
	var news model.ThemesOrLevelsNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := tc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Themes/Levels news created successfully",
		"data":    news,
	})
}

// UPDATE
func (tc *ThemesOrLevelsNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := tc.DB.Save(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.JSON(fiber.Map{
		"message": "Themes/Levels news updated successfully",
		"data":    news,
	})
}

// DELETE
// DELETE
func (tc *ThemesOrLevelsNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	if err := tc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Themes/Levels news with ID %v deleted successfully", news.ID),
	})
}

// Helper untuk update JSON field
func updateThemesOrLevelsNewsJSON(db *gorm.DB, themeID uint) {
	var newsList []model.ThemesOrLevelsNewsModel
	if err := db.Where("themes_or_level_id = ?", themeID).Order("created_at desc").Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch themes/levels news for update:", err)
		return
	}

	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal themes/levels news:", err)
		return
	}

	res := db.Table("themes_or_levels").
		Where("id = ?", themeID).
		Update("update_news", datatypes.JSON(newsData))

	log.Println("[DEBUG] Rows affected (themes_or_levels):", res.RowsAffected)
}
