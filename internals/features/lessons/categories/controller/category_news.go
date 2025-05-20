package controller

import (
	"fmt"
	"net/http"
	"quizku/internals/features/lessons/categories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryNewsController struct {
	DB *gorm.DB
}

func NewCategoryNewsController(db *gorm.DB) *CategoryNewsController {
	return &CategoryNewsController{DB: db}
}

// 🟢 GET ALL CATEGORY NEWS: Ambil semua data kategori berita
func (cc *CategoryNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.CategoryNewsModel

	if err := cc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil data category news",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Semua category news berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET CATEGORY NEWS BY CATEGORY_ID: Ambil semua berita berdasarkan category_id
func (cc *CategoryNewsController) GetByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	var news []model.CategoryNewsModel

	if err := cc.DB.
		Where("category_news_category_id = ?", categoryID).
		Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil category news berdasarkan category_id",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news berdasarkan category_id berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET CATEGORY NEWS BY ID: Ambil satu kategori berita berdasarkan ID
func (cc *CategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.CategoryNewsModel

	if err := cc.DB.
		Where("category_news_id = ?", id).
		First(&news).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news tidak ditemukan",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news berhasil ditemukan",
		"data":    news,
	})
}

// 🟢 CREATE CATEGORY NEWS: Tambahkan data kategori berita baru
func (cc *CategoryNewsController) Create(c *fiber.Ctx) error {
	var news model.CategoryNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Permintaan tidak valid (body tidak terbaca)",
		})
	}

	if err := cc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menyimpan category news",
			"detail":  err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Category news berhasil dibuat",
		"data":    news,
	})
}

// 🟢 UPDATE CATEGORY NEWS: Perbarui data kategori berita berdasarkan ID
func (cc *CategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var existing model.CategoryNewsModel

	if err := cc.DB.
		Where("category_news_id = ?", id).
		First(&existing).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news tidak ditemukan",
		})
	}

	var updateData model.CategoryNewsModel
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Permintaan tidak valid (body tidak terbaca)",
		})
	}

	updateData.CategoryNewsID = existing.CategoryNewsID // Pastikan ID tetap sama

	if err := cc.DB.Save(&updateData).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal memperbarui category news",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news berhasil diperbarui",
		"data":    updateData,
	})
}

// 🟢 DELETE CATEGORY NEWS: Hapus data kategori berita berdasarkan ID
func (cc *CategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.CategoryNewsModel

	if err := cc.DB.
		Where("category_news_id = ?", id).
		First(&news).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news tidak ditemukan",
		})
	}

	if err := cc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menghapus category news",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Category news dengan ID %v berhasil dihapus", news.CategoryNewsID),
	})
}
