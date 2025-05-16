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

	// 🔍 Query semua data dari database
	if err := sc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim data berita
	return c.JSON(fiber.Map{
		"message": "Subcategory news list retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET SUBCATEGORY NEWS BY SUBCATEGORY_ID: Ambil berita berdasarkan subcategory_id
func (sc *SubcategoryNewsController) GetBySubcategoryID(c *fiber.Ctx) error {
	subcategoryID := c.Params("subcategory_id")
	var news []model.SubcategoryNewsModel

	// 🔍 Query berita berdasarkan subcategory_id
	if err := sc.DB.Where("subcategory_id = ?", subcategoryID).Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve news by subcategory",
			"detail":  err.Error(),
		})
	}

	// ✅ Kirim data berita
	return c.JSON(fiber.Map{
		"message": "Subcategory news by subcategory retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET SUBCATEGORY NEWS BY ID: Ambil berita berdasarkan ID
func (sc *SubcategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	// 🔍 Cari berita berdasarkan ID
	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	// ✅ Kirim hasil pencarian
	return c.JSON(fiber.Map{
		"message": "Subcategory news found successfully",
		"data":    news,
	})
}

// 🟢 CREATE SUBCATEGORY NEWS: Tambahkan data berita subkategori baru
func (sc *SubcategoryNewsController) Create(c *fiber.Ctx) error {
	var news model.SubcategoryNewsModel

	// 🔄 Parsing body ke struct model
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// 💾 Simpan ke database
	if err := sc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Subcategory news created successfully",
		"data":    news,
	})
}

// 🟢 UPDATE SUBCATEGORY NEWS: Perbarui data berita berdasarkan ID
func (sc *SubcategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	// 🔍 Pastikan data ada
	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	// 🔄 Update dengan data baru dari body
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// 💾 Simpan update ke database
	if err := sc.DB.Save(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim hasil update
	return c.JSON(fiber.Map{
		"message": "Subcategory news updated successfully",
		"data":    news,
	})
}

// 🟢 DELETE SUBCATEGORY NEWS: Hapus berita berdasarkan ID
func (sc *SubcategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	// 🔍 Cek apakah berita ditemukan
	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	// 🗑️ Hapus dari database
	if err := sc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Konfirmasi penghapusan
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Subcategory news with ID %v deleted successfully", news.ID),
	})
}
