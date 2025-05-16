package controller

import (
	"fmt"
	"log"
	"quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ThemeOrLevelController menangani seluruh proses CRUD untuk themes atau levels (materi pelajaran)
type ThemeOrLevelController struct {
	DB *gorm.DB
}

// NewThemeOrLevelController membuat instance baru controller untuk themes_or_levels
func NewThemeOrLevelController(db *gorm.DB) *ThemeOrLevelController {
	return &ThemeOrLevelController{DB: db}
}

// 🟢 GET /themes-or-levels
// Mengambil semua data themes_or_levels yang tersedia di database
func (tc *ThemeOrLevelController) GetThemeOrLevels(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all themes or levels")
	var themesOrLevels []model.ThemesOrLevelsModel

	if err := tc.DB.Find(&themesOrLevels).Error; err != nil {
		log.Println("[ERROR] Failed to fetch themes or levels:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch themes or levels"})
	}

	log.Printf("[SUCCESS] Retrieved %d themes or levels\n", len(themesOrLevels))
	return c.JSON(fiber.Map{
		"message": "All themes or levels fetched successfully",
		"total":   len(themesOrLevels),
		"data":    themesOrLevels,
	})
}

// 🟢 GET /themes-or-levels/:id
// Mengambil satu data theme/level berdasarkan ID-nya
func (tc *ThemeOrLevelController) GetThemeOrLevelById(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching theme or level with ID:", id)

	var themeOrLevel model.ThemesOrLevelsModel
	if err := tc.DB.First(&themeOrLevel, id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Theme or level not found"})
	}

	log.Printf("[SUCCESS] Theme or level retrieved: ID=%d, Name=%s\n", themeOrLevel.ID, themeOrLevel.Name)
	return c.JSON(fiber.Map{
		"message": "Theme or level fetched successfully",
		"data":    themeOrLevel,
	})
}

// 🟢 GET /themes-or-levels/subcategory/:subcategory_id
// Mengambil semua themes/levels yang terkait dengan satu subkategori tertentu
func (tc *ThemeOrLevelController) GetThemesOrLevelsBySubcategory(c *fiber.Ctx) error {
	subcategoryID := c.Params("subcategory_id")
	log.Printf("[INFO] Fetching themes or levels for subcategory ID: %s\n", subcategoryID)

	var themesOrLevels []model.ThemesOrLevelsModel
	if err := tc.DB.Where("subcategories_id = ?", subcategoryID).Find(&themesOrLevels).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch themes or levels for subcategory ID %s: %v\n", subcategoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch themes or levels"})
	}

	log.Printf("[SUCCESS] Retrieved %d themes or levels for subcategory ID %s\n", len(themesOrLevels), subcategoryID)
	return c.JSON(fiber.Map{
		"message": "Themes or levels fetched successfully by subcategory",
		"total":   len(themesOrLevels),
		"data":    themesOrLevels,
	})
}

// 🟢 POST /themes-or-levels
// Menambahkan satu atau beberapa themes/levels sekaligus
func (tc *ThemeOrLevelController) CreateThemeOrLevel(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create theme or level")

	var single model.ThemesOrLevelsModel
	var multiple []model.ThemesOrLevelsModel

	// 🧠 Coba parse sebagai array terlebih dahulu
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d themes/levels as array\n", len(multiple))

		// ✅ Validasi tiap item array
		for i, item := range multiple {
			if item.Name == "" || item.Status == "" || item.DescriptionShort == "" || item.DescriptionLong == "" || item.SubcategoriesID == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "All fields are required in array (name, status, description_short, description_long, subcategories_id)",
					"index": i,
				})
			}
			if !isValidStatus(item.Status) {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid status in array. Allowed: active, pending, archived",
					"index": i,
				})
			}
			var count int64
			if err := tc.DB.Table("subcategories").Where("id = ?", item.SubcategoriesID).Count(&count).Error; err != nil || count == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid subcategories_id in array",
					"index": i,
				})
			}
		}

		if err := tc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple themes/levels: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create themes or levels"})
		}

		log.Printf("[SUCCESS] %d themes or levels created successfully\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple themes or levels created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Jika bukan array, parse sebagai objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single theme/level: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	log.Printf("[DEBUG] Parsed single theme/level: %+v\n", single)

	if single.Name == "" || single.Status == "" || single.DescriptionShort == "" || single.DescriptionLong == "" || single.SubcategoriesID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required (name, status, description_short, description_long, subcategories_id)"})
	}
	if !isValidStatus(single.Status) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status. Allowed: active, pending, archived"})
	}
	var count int64
	if err := tc.DB.Table("subcategories").Where("id = ?", single.SubcategoriesID).Count(&count).Error; err != nil || count == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid subcategories_id. The referenced subcategory does not exist"})
	}

	if err := tc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert single theme/level: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create theme or level"})
	}

	log.Printf("[SUCCESS] Theme or level created: ID=%d, Name=%s\n", single.ID, single.Name)
	return c.Status(201).JSON(fiber.Map{
		"message": "Theme or level created successfully",
		"data":    single,
	})
}

// isValidStatus memastikan status hanya berisi nilai yang diizinkan
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{"active": true, "pending": true, "archived": true}
	return validStatuses[status]
}

// 🟡 PUT /themes-or-levels/:id
// Memperbarui data theme/level berdasarkan ID
func (tc *ThemeOrLevelController) UpdateThemeOrLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating theme or level with ID:", id)

	var themeOrLevel model.ThemesOrLevelsModel
	if err := tc.DB.First(&themeOrLevel, id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Theme or level not found"})
	}

	if err := c.BodyParser(&themeOrLevel); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := tc.DB.Save(&themeOrLevel).Error; err != nil {
		log.Println("[ERROR] Failed to update theme or level:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update theme or level"})
	}

	log.Printf("[SUCCESS] Theme or level updated: ID=%d\n", themeOrLevel.ID)
	return c.JSON(fiber.Map{
		"message": "Theme or level updated successfully",
		"data":    themeOrLevel,
	})
}

// 🔴 DELETE /themes-or-levels/:id
// Menghapus data theme/level berdasarkan ID
func (tc *ThemeOrLevelController) DeleteThemeOrLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting theme or level with ID:", id)

	var theme model.ThemesOrLevelsModel
	if err := tc.DB.First(&theme, id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Theme or level not found",
		})
	}

	if err := tc.DB.Delete(&theme).Error; err != nil {
		log.Println("[ERROR] Failed to delete theme or level:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete theme or level",
		})
	}

	log.Printf("[SUCCESS] Theme or level with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Theme or level with ID %s deleted successfully", id),
	})
}
