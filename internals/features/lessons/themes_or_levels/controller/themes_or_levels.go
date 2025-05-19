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
	if err := tc.DB.First(&themeOrLevel, "themes_or_level_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Theme or level not found"})
	}

	log.Printf("[SUCCESS] Theme or level retrieved: ID=%d, Name=%s\n", themeOrLevel.ThemesOrLevelID, themeOrLevel.ThemesOrLevelName)
	return c.JSON(fiber.Map{
		"message": "Theme or level fetched successfully",
		"data":    themeOrLevel,
	})
}

// 🟢 GET /themes-or-levels/subcategory/:subcategory_id
// Mengambil semua themes/levels yang terkait dengan satu subkategori tertentu
func (tc *ThemeOrLevelController) GetThemesOrLevelsBySubcategory(c *fiber.Ctx) error {
	subcategoryID := c.Params("subcategory_id")
	log.Printf("[INFO] Fetching themes_or_levels for subcategory ID: %s\n", subcategoryID)

	var themesOrLevels []model.ThemesOrLevelsModel
	if err := tc.DB.
		Where("themes_or_level_subcategory_id = ?", subcategoryID).
		Find(&themesOrLevels).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch themes_or_levels for subcategory ID %s: %v\n", subcategoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil themes/levels"})
	}

	log.Printf("[SUCCESS] Retrieved %d themes_or_levels for subcategory ID %s\n", len(themesOrLevels), subcategoryID)
	return c.JSON(fiber.Map{
		"message": "Themes or levels fetched successfully by subcategory",
		"total":   len(themesOrLevels),
		"data":    themesOrLevels,
	})
}

// 🟢 POST /themes-or-levels
// Menambahkan satu atau beberapa themes/levels sekaligus
func (tc *ThemeOrLevelController) CreateThemeOrLevel(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create themes_or_levels")

	var single model.ThemesOrLevelsModel
	var multiple []model.ThemesOrLevelsModel

	// 🧠 Coba parse sebagai array terlebih dahulu
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d themes_or_levels as array\n", len(multiple))

		// ✅ Validasi tiap item array
		for i, item := range multiple {
			if item.ThemesOrLevelName == "" || item.ThemesOrLevelStatus == "" || item.ThemesOrLevelDescriptionShort == "" || item.ThemesOrLevelDescriptionLong == "" || item.ThemesOrLevelSubcategoryID == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "Semua field wajib diisi: themes_or_level_name, themes_or_level_status, themes_or_level_description_short, themes_or_level_description_long, themes_or_level_subcategory_id",
					"index": i,
				})
			}
			if !isValidStatus(item.ThemesOrLevelStatus) {
				return c.Status(400).JSON(fiber.Map{
					"error": "Status tidak valid. Hanya boleh: active, pending, archived",
					"index": i,
				})
			}
			var count int64
			if err := tc.DB.Table("subcategories").Where("id = ?", item.ThemesOrLevelSubcategoryID).Count(&count).Error; err != nil || count == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "themes_or_level_subcategory_id tidak valid",
					"index": i,
				})
			}
		}

		if err := tc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Gagal menyimpan multiple themes_or_levels: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat themes_or_levels"})
		}

		log.Printf("[SUCCESS] Berhasil menyimpan %d themes_or_levels\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Themes or levels created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Jika bukan array, parse sebagai objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Gagal parse single themes_or_level: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Body request tidak valid"})
	}
	log.Printf("[DEBUG] Parsed single theme_or_level: %+v\n", single)

	if single.ThemesOrLevelName == "" || single.ThemesOrLevelStatus == "" || single.ThemesOrLevelDescriptionShort == "" || single.ThemesOrLevelDescriptionLong == "" || single.ThemesOrLevelSubcategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Semua field wajib diisi: themes_or_level_name, themes_or_level_status, themes_or_level_description_short, themes_or_level_description_long, themes_or_level_subcategory_id"})
	}
	if !isValidStatus(single.ThemesOrLevelStatus) {
		return c.Status(400).JSON(fiber.Map{"error": "Status tidak valid. Hanya boleh: active, pending, archived"})
	}
	var count int64
	if err := tc.DB.Table("subcategories").Where("id = ?", single.ThemesOrLevelSubcategoryID).Count(&count).Error; err != nil || count == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "themes_or_level_subcategory_id tidak valid"})
	}

	if err := tc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Gagal menyimpan theme_or_level: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat theme_or_level"})
	}

	log.Printf("[SUCCESS] Theme_or_level created: ID=%d, Name=%s\n", single.ThemesOrLevelID, single.ThemesOrLevelName)
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

	var theme model.ThemesOrLevelsModel
	if err := tc.DB.First(&theme, "themes_or_level_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Theme or level not found"})
	}

	var updatedData model.ThemesOrLevelsModel
	if err := c.BodyParser(&updatedData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi data penting
	if updatedData.ThemesOrLevelName == "" || updatedData.ThemesOrLevelStatus == "" ||
		updatedData.ThemesOrLevelDescriptionShort == "" || updatedData.ThemesOrLevelDescriptionLong == "" ||
		updatedData.ThemesOrLevelSubcategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All required fields must be provided"})
	}
	if !isValidStatus(updatedData.ThemesOrLevelStatus) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status value"})
	}

	// Update field manual agar tidak menimpa kolom sensitif
	theme.ThemesOrLevelName = updatedData.ThemesOrLevelName
	theme.ThemesOrLevelStatus = updatedData.ThemesOrLevelStatus
	theme.ThemesOrLevelDescriptionShort = updatedData.ThemesOrLevelDescriptionShort
	theme.ThemesOrLevelDescriptionLong = updatedData.ThemesOrLevelDescriptionLong
	theme.ThemesOrLevelImageURL = updatedData.ThemesOrLevelImageURL
	theme.ThemesOrLevelSubcategoryID = updatedData.ThemesOrLevelSubcategoryID

	if err := tc.DB.Save(&theme).Error; err != nil {
		log.Println("[ERROR] Failed to update theme or level:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update theme or level"})
	}

	log.Printf("[SUCCESS] Theme or level updated: ID=%d\n", theme.ThemesOrLevelID)
	return c.JSON(fiber.Map{
		"message": "Theme or level updated successfully",
		"data":    theme,
	})
}

// 🔴 DELETE /themes-or-levels/:id
// Menghapus data theme/level berdasarkan ID
func (tc *ThemeOrLevelController) DeleteThemeOrLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting theme or level with ID:", id)

	var theme model.ThemesOrLevelsModel
	if err := tc.DB.First(&theme, "themes_or_level_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Theme or level not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Theme or level not found"})
	}

	if err := tc.DB.Delete(&theme).Error; err != nil {
		log.Println("[ERROR] Failed to delete theme or level:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete theme or level"})
	}

	log.Printf("[SUCCESS] Theme or level with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Theme or level with ID %s deleted successfully", id),
	})
}
