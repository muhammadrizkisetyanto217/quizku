package controller

import (
	"fmt"
	"net/http"
	"quizku/internals/features/lessons/subcategory/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubcategoryNewsController struct {
	DB *gorm.DB
}

func NewSubcategoryNewsController(db *gorm.DB) *SubcategoryNewsController {
	return &SubcategoryNewsController{DB: db}
}

// GET all subcategory news
func (sc *SubcategoryNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.SubcategoryNewsModel
	if err := sc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Subcategory news list retrieved successfully",
		"data":    news,
	})
}

func (sc *SubcategoryNewsController) GetBySubcategoryID(c *fiber.Ctx) error {
	subcategoryID := c.Params("subcategory_id")
	var news []model.SubcategoryNewsModel

	if err := sc.DB.Where("subcategory_id = ?", subcategoryID).Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve news by subcategory",
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subcategory news by subcategory retrieved successfully",
		"data":    news,
	})
}

// GET by ID
func (sc *SubcategoryNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subcategory news found successfully",
		"data":    news,
	})
}

// CREATE
func (sc *SubcategoryNewsController) Create(c *fiber.Ctx) error {
	var news model.SubcategoryNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := sc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Subcategory news created successfully",
		"data":    news,
	})
}

// UPDATE
func (sc *SubcategoryNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := sc.DB.Save(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subcategory news updated successfully",
		"data":    news,
	})
}

// DELETE
// DELETE
func (sc *SubcategoryNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.SubcategoryNewsModel

	if err := sc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Subcategory news not found",
		})
	}

	if err := sc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Subcategory news with ID %v deleted successfully", news.ID),
	})
}
