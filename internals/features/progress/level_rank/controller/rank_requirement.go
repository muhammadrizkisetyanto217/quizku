package controller

import (
	"quizku/internals/features/progress/level_rank/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RankRequirementController struct {
	DB *gorm.DB
}

func NewRankRequirementController(db *gorm.DB) *RankRequirementController {
	return &RankRequirementController{DB: db}
}

// CREATE MANY
func (ctrl *RankRequirementController) Create(c *fiber.Ctx) error {
	var inputs []model.RankRequirement
	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak valid"})
	}
	if len(inputs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak boleh kosong"})
	}
	if err := ctrl.DB.Create(&inputs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Data berhasil ditambahkan",
		"data":    inputs,
	})
}

// GET ALL
func (ctrl *RankRequirementController) GetAll(c *fiber.Ctx) error {
	var ranks []model.RankRequirement
	if err := ctrl.DB.Order("rank asc").Find(&ranks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ranks)
}

// GET BY ID
func (ctrl *RankRequirementController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var rank model.RankRequirement
	if err := ctrl.DB.First(&rank, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	return c.JSON(rank)
}

// UPDATE BY ID
func (ctrl *RankRequirementController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var input model.RankRequirement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tidak valid"})
	}
	var existing model.RankRequirement
	if err := ctrl.DB.First(&existing, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	input.ID = existing.ID
	if err := ctrl.DB.Save(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Data berhasil diupdate", "data": input})
}

// DELETE BY ID
func (ctrl *RankRequirementController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctrl.DB.Delete(&model.RankRequirement{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Data berhasil dihapus"})
}
