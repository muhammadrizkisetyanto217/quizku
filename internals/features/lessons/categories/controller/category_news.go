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
	var categories []model.CategoryNewsModel

	// 🔍 Query semua data dari database
	if err := cc.DB.Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.JSON(fiber.Map{
		"message": "All category news retrieved successfully",
		"data":    categories,
	})
}

// 🟢 GET CATEGORY NEWS BY CATEGORY_ID: Ambil semua berita berdasarkan category_id
func (cc *CategoryNewsController) GetByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	var categories []model.CategoryNewsModel
	if err := cc.DB.Where("category_id = ?", categoryID).Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.JSON(fiber.Map{
		"message": "Category news filtered by category_id retrieved successfully",
		"data":    categories,
	})
}

// 🟢 GET CATEGORY NEWS BY ID: Ambil satu kategori berita berdasarkan ID
func (cc *CategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	// 🔍 Cari data berdasarkan ID
	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	// ✅ Kirim respons sukses
	return c.JSON(fiber.Map{
		"message": "Category news found successfully",
		"data":    category,
	})
}

// 🟢 CREATE CATEGORY NEWS: Tambahkan data kategori berita baru
func (cc *CategoryNewsController) Create(c *fiber.Ctx) error {
	var category model.CategoryNewsModel

	// 🔄 Parsing body request
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// 💾 Simpan ke database
	if err := cc.DB.Create(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Category news created successfully",
		"data":    category,
	})
}

// 🟢 UPDATE CATEGORY NEWS: Perbarui data kategori berita berdasarkan ID
func (cc *CategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	// 🔍 Cek apakah data dengan ID tersebut ada
	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	// 🔄 Parsing body baru ke struct
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// 💾 Simpan perubahan
	if err := cc.DB.Save(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.JSON(fiber.Map{
		"message": "Category news updated successfully",
		"data":    category,
	})
}

// 🟢 DELETE CATEGORY NEWS: Hapus data kategori berita berdasarkan ID
func (cc *CategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	// 🔍 Cek apakah data ada
	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	// 🗑️ Hapus data dari database
	if err := cc.DB.Delete(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// ✅ Kirim respons sukses
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Category news with ID %v deleted successfully", category.ID),
	})
}
