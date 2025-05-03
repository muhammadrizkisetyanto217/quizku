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
		"message": "All category news retrieved successfully",
		"data":    categories,
	})
}

// GET by category_id
func (cc *CategoryNewsController) GetByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	var categories []model.CategoryNewsModel
	if err := cc.DB.Where("category_id = ?", categoryID).Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category news filtered by category_id retrieved successfully",
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

	return c.JSON(fiber.Map{
		"message": "Category news updated successfully",
		"data":    category,
	})
}

// DELETE
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

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Category news with ID %v deleted successfully", category.ID),
	})
}
