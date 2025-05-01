package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"quizku/internals/features/lessons/categories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryNewsController struct {
	DB *gorm.DB
}

func NewCategoryNewsController(db *gorm.DB) *CategoryNewsController {
	return &CategoryNewsController{DB: db}
}

// GET all category news
func (cc *CategoryNewsController) GetAll(c *fiber.Ctx) error {
	var categories []model.CategoryNewsModel
	if err := cc.DB.Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news list retrieved successfully",
		"data":    categories,
	})
}

// GET by ID
func (cc *CategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news found successfully",
		"data":    category,
	})
}

// CREATE
func (cc *CategoryNewsController) Create(c *fiber.Ctx) error {
	var category model.CategoryNewsModel

	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := cc.DB.Create(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateCategoryNewsJSON(cc.DB, category.CategoryID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Category news created successfully",
		"data":    category,
	})
}

// UPDATE
func (cc *CategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := cc.DB.Save(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateCategoryNewsJSON(cc.DB, category.CategoryID)

	return c.JSON(fiber.Map{
		"message": "Category news updated successfully",
		"data":    category,
	})
}

// DELETE
func (cc *CategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.CategoryNewsModel

	if err := cc.DB.First(&category, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Category news not found",
		})
	}

	if err := cc.DB.Delete(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateCategoryNewsJSON(cc.DB, category.CategoryID)

	return c.JSON(fiber.Map{
		"message": "Category news deleted successfully",
	})
}

// Helper untuk update kolom update_news di tabel categories
func updateCategoryNewsJSON(db *gorm.DB, categoryID int) {
	var newsList []model.CategoryNewsModel
	if err := db.Where("category_id = ?", categoryID).Order("created_at desc").Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch category news for update:", err)
		return
	}

	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal category news:", err)
		return
	}

	res := db.Table("categories").
		Where("id = ?", categoryID).
		Update("update_news", datatypes.JSON(newsData))

	log.Println("[DEBUG] Rows affected:", res.RowsAffected)
}
