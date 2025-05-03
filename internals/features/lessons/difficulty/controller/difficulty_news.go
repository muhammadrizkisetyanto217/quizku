package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"quizku/internals/features/lessons/difficulty/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DifficultyNewsController struct {
	DB *gorm.DB
}

func NewDifficultyNewsController(db *gorm.DB) *DifficultyNewsController {
	return &DifficultyNewsController{DB: db}
}

// GET all news
func (dc *DifficultyNewsController) GetAllNews(c *fiber.Ctx) error {
	var newsList []model.DifficultyNews

	if err := dc.DB.Order("created_at desc").Find(&newsList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil semua berita",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Semua news berhasil diambil",
		"data":    newsList,
	})
}

// GET all news by difficulty
func (dc *DifficultyNewsController) GetNewsByDifficulty(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Println("[DEBUG] Difficulty ID:", difficultyID)

	var news []model.DifficultyNews
	if err := dc.DB.Where("difficulty_id = ?", difficultyID).Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	log.Printf("[DEBUG] Fetched %d news\n", len(news))

	return c.JSON(fiber.Map{
		"message": "News list retrieved successfully",
		"data":    news,
	})
}

// GET news by ID
func (dc *DifficultyNewsController) GetNewsByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNews

	if err := dc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News found successfully",
		"data":    news,
	})
}

// CREATE news
func (dc *DifficultyNewsController) CreateNews(c *fiber.Ctx) error {
	var news model.DifficultyNews

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := dc.DB.Create(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateNewsJSON(dc.DB, news.DifficultyID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "News created successfully",
		"data":    news,
	})
}

// UPDATE news
func (dc *DifficultyNewsController) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNews

	if err := dc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News not found",
		})
	}

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := dc.DB.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateNewsJSON(dc.DB, news.DifficultyID)

	return c.JSON(fiber.Map{
		"message": "News updated successfully",
		"data":    news,
	})
}

// DELETE news
func (dc *DifficultyNewsController) DeleteNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNews

	if err := dc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News not found",
		})
	}

	if err := dc.DB.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateNewsJSON(dc.DB, news.DifficultyID)

	return c.JSON(fiber.Map{
		"message": "News deleted successfully",
	})
}

// Helper untuk update kolom update_news di difficulties
func updateNewsJSON(db *gorm.DB, difficultyID uint) {
	var newsList []model.DifficultyNews
	if err := db.Where("difficulty_id = ?", difficultyID).Order("created_at desc").Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch news for update:", err)
		return
	}

	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal news:", err)
		return
	}

	db.Model(&model.DifficultyModel{}).
		Where("id = ?", difficultyID).
		Update("update_news", datatypes.JSON(newsData))
}
