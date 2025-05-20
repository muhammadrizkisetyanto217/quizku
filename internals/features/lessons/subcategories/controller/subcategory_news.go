package controller

import (
	"fmt"
	"net/http"
	"quizku/internals/features/lessons/subcategories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubcategoryNewsController struct {
	DB *gorm.DB
}

func NewSubcategoryNewsController(db *gorm.DB) *SubcategoryNewsController {
	return &SubcategoryNewsController{DB: db}
}

// 🟢 GET ALL SUBCATEGORY NEWS: Ambil seluruh data berita subkategori
func (sc *SubcategoryNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.SubcategoryNewsModel

	if err := sc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil data subcategory news",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita subkategori berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET SUBCATEGORY NEWS BY SUBCATEGORY_ID: Ambil berita berdasarkan subcategory
func (sc *SubcategoryNewsController) GetBySubcategoryID(c *fiber.Ctx) error {
	subcategoryID := c.Params("subcategory_id")
	var news []model.SubcategoryNewsModel

	if err := sc.DB.
		Where("subcategory_news_subcategory_id = ?", subcategoryID).
		Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil berita berdasarkan subcategory_id",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita subkategori berdasarkan subcategory_id berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET SUBCATEGORY NEWS BY ID: Ambil satu berita berdasarkan ID
func (sc *SubcategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	if err := sc.DB.
		Where("subcategory_news_id = ?", id).
		First(&news).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita subkategori tidak ditemukan",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita subkategori berhasil ditemukan",
		"data":    news,
	})
}

// 🟢 CREATE SUBCATEGORY NEWS: Tambahkan data berita subkategori baru
func (sc *SubcategoryNewsController) Create(c *fiber.Ctx) error {
	var news model.SubcategoryNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Permintaan tidak valid (body tidak terbaca)",
		})
	}

	if err := sc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menyimpan berita",
			"detail":  err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Berita subkategori berhasil dibuat",
		"data":    news,
	})
}

// 🟢 UPDATE SUBCATEGORY NEWS: Perbarui data berita berdasarkan ID
func (sc *SubcategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var existing model.SubcategoryNewsModel

	// 🔍 Cek apakah berita dengan ID tersebut ada
	if err := sc.DB.Where("subcategory_news_id = ?", id).First(&existing).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita subkategori tidak ditemukan",
		})
	}

	// 🔄 Update dengan body baru
	var updateData model.SubcategoryNewsModel
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Permintaan tidak valid (body tidak terbaca)",
		})
	}

	// Tetap pakai ID yang lama
	updateData.SubcategoryNewsID = existing.SubcategoryNewsID

	if err := sc.DB.Save(&updateData).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal memperbarui berita",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita subkategori berhasil diperbarui",
		"data":    updateData,
	})
}

// 🟢 DELETE SUBCATEGORY NEWS: Hapus berita berdasarkan ID
func (sc *SubcategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	if err := sc.DB.Where("subcategory_news_id = ?", id).First(&news).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita subkategori tidak ditemukan",
		})
	}

	if err := sc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menghapus berita",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Berita subkategori dengan ID %v berhasil dihapus", news.SubcategoryNewsID),
	})
}
