package controller

import (
	"log"
	"quizku/internals/features/donations/donation_questions/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DonationQuestionController struct {
	DB *gorm.DB
}

func NewDonationQuestionController(db *gorm.DB) *DonationQuestionController {
	return &DonationQuestionController{DB: db}
}

// GET all donation_questions
func (ctrl *DonationQuestionController) GetAll(c *fiber.Ctx) error {
	var items []model.DonationQuestionModel
	if err := ctrl.DB.Find(&items).Error; err != nil {
		log.Println("[ERROR] Failed to fetch donation questions:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch donation questions"})
	}

	return c.JSON(fiber.Map{"data": items})
}

// GET by ID
func (ctrl *DonationQuestionController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item model.DonationQuestionModel

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		log.Println("[ERROR] Donation question not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Donation question not found"})
	}

	return c.JSON(item)
}

func (ctrl *DonationQuestionController) GetByDonationID(c *fiber.Ctx) error {
	donationID := c.Params("donationId")
	var items []model.DonationQuestionModel

	if err := ctrl.DB.
		Where("donation_id = ?", donationID).
		Find(&items).Error; err != nil {
		log.Println("[ERROR] Failed to fetch donation questions by donation_id:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch by donation ID"})
	}

	return c.JSON(fiber.Map{"data": items})
}

// POST create new donation_question
func (ctrl *DonationQuestionController) Create(c *fiber.Ctx) error {
	var input model.DonationQuestionModel
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create donation question:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create donation question"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Donation question created successfully",
		"data":    input,
	})
}

// PUT update donation_question
func (ctrl *DonationQuestionController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var item model.DonationQuestionModel

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Donation question not found"})
	}

	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := ctrl.DB.Save(&item).Error; err != nil {
		log.Println("[ERROR] Failed to update donation question:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update donation question"})
	}

	return c.JSON(fiber.Map{"message": "Updated successfully", "data": item})
}

// DELETE donation_question
func (ctrl *DonationQuestionController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.DonationQuestionModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete donation question:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete donation question"})
	}

	return c.JSON(fiber.Map{"message": "Donation question deleted successfully"})
}
