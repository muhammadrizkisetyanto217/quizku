package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"quizku/internals/features/progress/level_rank/model"
)

type LevelRequirementController struct {
	DB *gorm.DB
}

func NewLevelRequirementController(db *gorm.DB) *LevelRequirementController {
	return &LevelRequirementController{DB: db}
}

// GET /api/level-requirements
func (ctrl *LevelRequirementController) GetAll(c *fiber.Ctx) error {
	var levels []model.LevelRequirement
	if err := ctrl.DB.Order("level ASC").Find(&levels).Error; err != nil {
		log.Println("[ERROR] Gagal ambil level:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil level"})
	}
	return c.JSON(fiber.Map{"data": levels})
}

// GET /api/level-requirements/:id
func (ctrl *LevelRequirementController) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var level model.LevelRequirement
	if err := ctrl.DB.First(&level, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Level tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"data": level})
}


// POST /api/level-requirements
func (ctrl *LevelRequirementController) Create(c *fiber.Ctx) error {
	var input []model.LevelRequirement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON harus berupa array"})
	}

	if len(input) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data level kosong"})
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Gagal buat level batch:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Level berhasil ditambahkan",
		"data":    input,
	})
}


// PUT /api/level-requirements/:id
func (ctrl *LevelRequirementController) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var level model.LevelRequirement
	if err := ctrl.DB.First(&level, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Level tidak ditemukan"})
	}
	var input model.LevelRequirement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	input.ID = level.ID // pastikan tidak berubah
	if err := ctrl.DB.Save(&input).Error; err != nil {
		log.Println("[ERROR] Gagal update level:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update level"})
	}
	return c.JSON(fiber.Map{"data": input})
}

// DELETE /api/level-requirements/:id
func (ctrl *LevelRequirementController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := ctrl.DB.Delete(&model.LevelRequirement{}, id).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal hapus level"})
	}
	return c.JSON(fiber.Map{"message": "Level berhasil dihapus"})
}
