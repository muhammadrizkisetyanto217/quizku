package controller

import (
	"fmt"
	"log"
	"net/http"
	"quizku/internals/features/lessons/difficulty/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DifficultyNewsController struct {
	DB *gorm.DB
}

func NewDifficultyNewsController(db *gorm.DB) *DifficultyNewsController {
	return &DifficultyNewsController{DB: db}
}

// GET all news
// 🟢 GET ALL DIFFICULTY NEWS: Ambil semua berita berdasarkan difficulty
func (dc *DifficultyNewsController) GetAllNews(c *fiber.Ctx) error {
	var newsList []model.DifficultyNewsModel

	if err := dc.DB.Order("difficulty_news_created_at DESC").Find(&newsList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil semua berita",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil semua berita difficulty",
		"data":    newsList,
	})
}

// GET all news by difficulty
// 🟢 GET DIFFICULTY NEWS BY DIFFICULTY ID
func (dc *DifficultyNewsController) GetNewsByDifficultyId(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Println("[DEBUG] difficulty_id:", difficultyID)

	var news []model.DifficultyNewsModel
	if err := dc.DB.
		Where("difficulty_news_difficulty_id = ?", difficultyID).
		Order("difficulty_news_created_at DESC").
		Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil news berdasarkan difficulty_id",
		})
	}

	log.Printf("[SUCCESS] Berhasil mengambil %d news untuk difficulty_id %s\n", len(news), difficultyID)

	return c.JSON(fiber.Map{
		"message": "Berita berdasarkan difficulty berhasil diambil",
		"data":    news,
	})
}

// GET news by ID
// 🟢 GET DIFFICULTY NEWS BY ID
func (dc *DifficultyNewsController) GetNewsByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNewsModel

	// 🔍 Ambil berdasarkan primary key
	if err := dc.DB.Where("difficulty_news_id = ?", id).First(&news).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita ditemukan",
		"data":    news,
	})
}

// CREATE news
// 🟢 CREATE DIFFICULTY NEWS
func (dc *DifficultyNewsController) CreateNews(c *fiber.Ctx) error {
	var news model.DifficultyNewsModel

	// 🔄 Parse body
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Format request tidak valid",
		})
	}

	// 💾 Simpan ke DB
	if err := dc.DB.Create(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menyimpan berita",
			"detail":  err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Berita berhasil dibuat",
		"data":    news,
	})
}


// UPDATE news
// 🟢 UPDATE DIFFICULTY NEWS
func (dc *DifficultyNewsController) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNewsModel

	// 🔍 Cari berita berdasarkan ID
	if err := dc.DB.Where("difficulty_news_id = ?", id).First(&news).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita tidak ditemukan",
		})
	}

	// 🔄 Parsing data baru dari body
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Request body tidak valid",
		})
	}

	// 💾 Simpan ke DB
	if err := dc.DB.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal memperbarui berita",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita berhasil diperbarui",
		"data":    news,
	})
}

// 🟢 DELETE DIFFICULTY NEWS
func (dc *DifficultyNewsController) DeleteNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.DifficultyNewsModel

	// 🔍 Cari berita berdasarkan ID
	if err := dc.DB.Where("difficulty_news_id = ?", id).First(&news).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita tidak ditemukan",
		})
	}

	// 🗑️ Soft delete
	if err := dc.DB.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menghapus berita",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Berita dengan ID %v berhasil dihapus", news.DifficultyNewsID),
	})
}
