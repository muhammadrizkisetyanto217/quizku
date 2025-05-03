package controller

import (
	"fmt"
	"log"
	"quizku/internals/features/lessons/units/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UnitController struct {
	DB *gorm.DB
}

func NewUnitController(db *gorm.DB) *UnitController {
	return &UnitController{DB: db}
}

// GET all units
func (uc *UnitController) GetUnits(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all units")
	var units []model.UnitModel

	if err := uc.DB.
		Preload("SectionQuizzes").
		Preload("SectionQuizzes.Quizzes"). // ✅ Preload ke quizzes
		Find(&units).Error; err != nil {

		log.Println("[ERROR] Failed to fetch units:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch units",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d units\n", len(units))
	return c.JSON(fiber.Map{
		"message": "All units fetched successfully",
		"total":   len(units),
		"data":    units,
	})
}

// GET single unit by ID
func (uc *UnitController) GetUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching unit with ID:", id)

	var unit model.UnitModel
	if err := uc.DB.Preload("SectionQuizzes").First(&unit, id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Unit not found"})
	}

	log.Printf("[SUCCESS] Unit found: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit fetched successfully",
		"data":    unit,
	})
}

// GET units by themes_or_level_id
func (uc *UnitController) GetUnitByThemesOrLevels(c *fiber.Ctx) error {
	themesOrLevelID := c.Params("themesOrLevelId")
	log.Printf("[INFO] Fetching units with themes_or_level_id: %s\n", themesOrLevelID)

	var units []model.UnitModel
	if err := uc.DB.Preload("SectionQuizzes").
		Where("themes_or_level_id = ?", themesOrLevelID).
		Find(&units).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch units for themes_or_level_id %s: %v\n", themesOrLevelID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch units"})
	}

	log.Printf("[SUCCESS] Retrieved %d units for themes_or_level_id %s\n", len(units), themesOrLevelID)
	return c.JSON(fiber.Map{
		"message": "Units fetched successfully by themes_or_level",
		"total":   len(units),
		"data":    units,
	})
}

// CreateUnit menangani input satu atau banyak unit
func (uc *UnitController) CreateUnit(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create unit")

	var (
		single   model.UnitModel
		multiple []model.UnitModel
	)

	raw := c.Body() // Ambil raw body
	if len(raw) > 0 && raw[0] == '[' {
		// Dipastikan ini array
		if err := c.BodyParser(&multiple); err != nil {
			log.Printf("[ERROR] Failed to parse unit array: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON array"})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Received empty array of units")
			return c.Status(400).JSON(fiber.Map{"error": "Array of units is empty"})
		}

		// Validasi setiap item
		for i, unit := range multiple {
			if unit.ThemesOrLevelID == 0 || unit.Name == "" {
				log.Printf("[ERROR] Invalid unit at index %d: %+v\n", i, unit)
				return c.Status(400).JSON(fiber.Map{
					"error":      "Each unit must have a valid themes_or_level_id and name",
					"index":      i,
					"unit_input": unit,
				})
			}
		}

		// Insert batch
		if err := uc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple units: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create units"})
		}

		log.Printf("[SUCCESS] Inserted %d units", len(multiple))
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple units created successfully",
			"data":    multiple,
		})
	}

	// Fallback: parse single object
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single unit input: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format (expected object or array)",
		})
	}

	log.Printf("[DEBUG] Parsed single unit: %+v\n", single)

	if single.ThemesOrLevelID == 0 || single.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "themes_or_level_id and name are required"})
	}

	if err := uc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert unit: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create unit"})
	}

	log.Printf("[SUCCESS] Unit created: ID=%d, Name=%s\n", single.ID, single.Name)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Unit created successfully",
		"data":    single,
	})
}

// UPDATE unit
func (uc *UnitController) UpdateUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating unit with ID:", id)

	var unit model.UnitModel
	if err := uc.DB.First(&unit, id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Unit not found"})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := uc.DB.Model(&unit).Updates(requestData).Error; err != nil {
		log.Println("[ERROR] Failed to update unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update unit"})
	}

	log.Printf("[SUCCESS] Unit updated: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit updated successfully",
		"data":    unit,
	})
}

// DELETE unit
func (uc *UnitController) DeleteUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting unit with ID:", id)

	if err := uc.DB.Delete(&model.UnitModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete unit",
		})
	}

	log.Printf("[SUCCESS] Unit with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Unit with ID %s deleted successfully", id),
	})
}
