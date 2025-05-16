package controller

import (
	"encoding/json"
	"fmt"
	"log"

	dto "quizku/internals/features/lessons/categories/dto"
	"quizku/internals/features/lessons/categories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

// 🟢 GET ALL CATEGORIES: Ambil semua kategori dari database
func (cc *CategoryController) GetCategories(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all categories")

	// 🔍 Ambil semua data kategori dari database
	var categories []model.CategoryModel
	if err := cc.DB.Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	// ✅ Log & kirim data jika berhasil
	log.Printf("[SUCCESS] Retrieved %d categories\n", len(categories))
	return c.JSON(fiber.Map{
		"message": "All categories fetched successfully",
		"total":   len(categories),
		"data":    categories,
	})
}

// 🟢 GET CATEGORY BY ID: Ambil satu kategori berdasarkan ID
func (cc *CategoryController) GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching category with ID: %s\n", id)

	var category model.CategoryModel

	// 🔍 Cari data berdasarkan ID
	if err := cc.DB.First(&category, id).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// ✅ Berhasil ditemukan
	log.Printf("[SUCCESS] Retrieved category: ID=%s, Name=%s\n", id, category.Name)
	return c.JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data":    category,
	})
}

// 🟢 GET CATEGORY NAMES BY DIFFICULTY ID: Ambil nama kategori berdasarkan difficulty_id
func (cc *CategoryController) GetCategoriesByDifficulty(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Printf("[INFO] Fetching categories with difficulty ID: %s\n", difficultyID)

	var categories []model.CategoryModel

	// 🔍 Ambil hanya ID dan Name berdasarkan difficulty_id
	if err := cc.DB.
		Select("id", "name").
		Where("difficulty_id = ?", difficultyID).
		Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories for difficulty ID %s: %v\n", difficultyID, err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	// 🔄 Format response ke dalam DTO (tooltip version)
	var responses []dto.CategoryTooltipResponse
	for _, c := range categories {
		responses = append(responses, dto.CategoryTooltipResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	// ✅ Berhasil ambil data
	log.Printf("[SUCCESS] Retrieved %d category tooltips for difficulty ID %s\n", len(responses), difficultyID)
	return c.JSON(fiber.Map{
		"message": "Category names fetched successfully",
		"total":   len(responses),
		"data":    responses,
	})
}

// 🟢 CREATE CATEGORY: Tambahkan satu atau banyak kategori
func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create category")

	// 🔄 Coba parsing array kategori terlebih dahulu
	var singleCategory model.CategoryModel
	var multipleCategories []model.CategoryModel

	if err := c.BodyParser(&multipleCategories); err == nil && len(multipleCategories) > 0 {
		// 💾 Simpan multiple kategori
		if err := cc.DB.Create(&multipleCategories).Error; err != nil {
			log.Printf("[ERROR] Failed to create multiple categories: %v\n", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create multiple categories",
			})
		}
		log.Printf("[SUCCESS] %d categories created\n", len(multipleCategories))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple categories created successfully",
			"data":    multipleCategories,
		})
	}

	// 🔄 Jika bukan array, coba parsing objek tunggal
	if err := c.BodyParser(&singleCategory); err != nil {
		log.Printf("[ERROR] Invalid input for single category: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// 💾 Simpan single kategori
	if err := cc.DB.Create(&singleCategory).Error; err != nil {
		log.Printf("[ERROR] Failed to create single category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	log.Printf("[SUCCESS] Category created: ID=%d, Name=%s\n", singleCategory.ID, singleCategory.Name)
	return c.Status(201).JSON(fiber.Map{
		"message": "Category created successfully",
		"data":    singleCategory,
	})
}

// 🟢 UPDATE CATEGORY: Perbarui kategori berdasarkan ID
func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating category with ID: %s\n", id)

	// 🔍 Cari data berdasarkan ID
	var category model.CategoryModel
	if err := cc.DB.First(&category, id).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// 🔄 Parsing input ke dalam map
	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Invalid input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// 🧩 Khusus untuk field update_news (type JSON)
	if un, ok := input["update_news"]; ok {
		jsonData, err := json.Marshal(un)
		if err == nil {
			input["update_news"] = datatypes.JSON(jsonData)
		}
	}

	// 💾 Lakukan update
	if err := cc.DB.Model(&category).Updates(input).Error; err != nil {
		log.Printf("[ERROR] Failed to update category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	log.Printf("[SUCCESS] Category updated: ID=%s, Name=%s\n", id, category.Name)
	return c.JSON(fiber.Map{
		"message": "Category updated successfully",
		"data":    category,
	})
}

// 🟢 DELETE CATEGORY: Hapus kategori berdasarkan ID
func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting category with ID: %s\n", id)

	// 🗑️ Hapus langsung berdasarkan ID tanpa preload
	if err := cc.DB.Delete(&model.CategoryModel{}, id).Error; err != nil {
		log.Printf("[ERROR] Failed to delete category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	log.Printf("[SUCCESS] Category with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Category with ID %s deleted successfully", id),
	})
}
