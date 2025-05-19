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

// 🟢 GET /api/units
// Mengambil semua unit yang tersedia di database, lengkap dengan relasi ke SectionQuizzes dan Quizzes.
// Biasanya digunakan untuk halaman admin, pembelajaran, atau struktur modul.
func (uc *UnitController) GetUnits(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all units")
	var units []model.UnitModel

	// Preload relasi ke section_quizzes dan quizzes di dalamnya
	if err := uc.DB.
		Preload("SectionQuizzes").
		Preload("SectionQuizzes.Quizzes"). // ✅ Preload nested relasi
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

// 🟢 GET /api/units/:id
// Mengambil satu unit berdasarkan ID-nya, lengkap dengan section_quizzes.
// Cocok digunakan saat membuka halaman detail unit.
func (uc *UnitController) GetUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching unit with ID:", id)

	var unit model.UnitModel

	// Preload SectionQuizzes saja (tidak sampai ke quizzes)
	if err := uc.DB.Preload("SectionQuizzes").First(&unit, "unit_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Unit not found"})
	}

	log.Printf("[SUCCESS] Unit found: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit fetched successfully",
		"data":    unit,
	})
}

// 🟢 GET /api/units/themes-or-levels/:themesOrLevelId
// Mengambil semua unit berdasarkan themes_or_level_id.
// Cocok untuk menampilkan semua unit di dalam 1 tema atau level materi.
func (uc *UnitController) GetUnitByThemesOrLevels(c *fiber.Ctx) error {
	themesOrLevelID := c.Params("themesOrLevelId")
	log.Printf("[INFO] Fetching units with themes_or_level_id: %s\n", themesOrLevelID)

	var units []model.UnitModel

	// Ambil semua unit berdasarkan foreign key themes_or_level_id
	if err := uc.DB.Preload("SectionQuizzes").
		Where("unit_themes_or_level_id = ?", themesOrLevelID).
		Find(&units).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch units for unit_themes_or_level_id %s: %v\n", themesOrLevelID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch units"})
	}

	log.Printf("[SUCCESS] Retrieved %d units for themes_or_level_id %s\n", len(units), themesOrLevelID)
	return c.JSON(fiber.Map{
		"message": "Units fetched successfully by themes_or_level",
		"total":   len(units),
		"data":    units,
	})
}

// 🟡 POST /api/units
// Membuat satu atau banyak unit sekaligus.
// - Jika body berisi array JSON → batch insert (banyak unit).
// - Jika body berisi objek JSON tunggal → insert satu unit.
// Field wajib: themes_or_level_id dan name.
func (uc *UnitController) CreateUnit(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create unit")

	var (
		single   model.UnitModel
		multiple []model.UnitModel
	)

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		if err := c.BodyParser(&multiple); err != nil {
			log.Printf("[ERROR] Failed to parse unit array: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON array"})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Received empty array of units")
			return c.Status(400).JSON(fiber.Map{"error": "Array of units is empty"})
		}

		for i, unit := range multiple {
			if unit.UnitThemesOrLevelID == 0 || unit.UnitName == "" {
				log.Printf("[ERROR] Invalid unit at index %d: %+v\n", i, unit)
				return c.Status(400).JSON(fiber.Map{
					"error":      "Each unit must have a valid unit_themes_or_level_id and unit_name",
					"index":      i,
					"unit_input": unit,
				})
			}
		}

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

	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single unit input: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format (expected object or array)",
		})
	}

	log.Printf("[DEBUG] Parsed single unit: %+v\n", single)

	if single.UnitThemesOrLevelID == 0 || single.UnitName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "unit_themes_or_level_id and unit_name are required"})
	}

	if err := uc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert unit: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create unit"})
	}

	log.Printf("[SUCCESS] Unit created: ID=%d, Name=%s\n", single.UnitID, single.UnitName)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Unit created successfully",
		"data":    single,
	})
}

// 🟠 PUT /api/units/:id
// Mengupdate unit berdasarkan ID.
// Field yang diupdate fleksibel karena menerima map[string]interface{} dari body.
// 🟠 PUT /api/units/:id
// Mengupdate unit berdasarkan unit_id.
// Field yang diupdate fleksibel karena menerima map[string]interface{} dari body.
func (uc *UnitController) UpdateUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating unit with unit_id:", id)

	var unit model.UnitModel

	// Cek apakah unit dengan unit_id tersebut ada
	if err := uc.DB.First(&unit, "unit_id = ?", id).Error; err != nil {
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

	log.Printf("[SUCCESS] Unit updated: unit_id=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit updated successfully",
		"data":    unit,
	})
}

// 🔴 DELETE /api/units/:id
// Menghapus unit berdasarkan ID.
// Menggunakan soft delete jika model menggunakan gorm.Model dengan DeletedAt.
func (uc *UnitController) DeleteUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting unit with unit_id:", id)

	var unit model.UnitModel

	if err := uc.DB.First(&unit, "unit_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Unit not found",
		})
	}

	if err := uc.DB.Delete(&unit).Error; err != nil {
		log.Println("[ERROR] Failed to delete unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete unit",
		})
	}

	log.Printf("[SUCCESS] Unit with unit_id %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Unit with unit_id %s deleted successfully", id),
	})
}
