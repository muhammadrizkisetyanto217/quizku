package utils

import (
	"fmt"
	"log"

	database "quizku/internals/databases"
	"quizku/internals/features/utils/tooltips/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// TooltipsController menangani semua operasi terkait tooltips
type TooltipsController struct {
	DB *gorm.DB
}

// NewTooltipsController membuat instance baru dari TooltipsController
func NewTooltipsController(db *gorm.DB) *TooltipsController {
	return &TooltipsController{DB: db}
}

func (tc *TooltipsController) GetTooltipsID(c *fiber.Ctx) error {
	log.Println("Fetching tooltips for given keywords")

	var request struct {
		Keywords []string `json:"keywords"`
	}

	// Parsing request body
	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Inisialisasi array untuk menyimpan ID tooltips
	var tooltipIDs []uint

	// Loop pencarian keyword dalam database
	for _, keyword := range request.Keywords {
		var tooltip model.Tooltip
		if err := database.DB.Select("id").Where("keyword = ?", keyword).First(&tooltip).Error; err == nil {
			tooltipIDs = append(tooltipIDs, tooltip.ID)
		}
	}

	// Mengembalikan array ID tooltips
	return c.JSON(fiber.Map{
		"tooltips_id": tooltipIDs,
	})
}

// InsertTooltip menangani permintaan untuk menambahkan tooltips baru
func (tc *TooltipsController) CreateTooltip(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create tooltip(s)")

	var (
		single   model.Tooltip
		multiple []model.Tooltip
	)

	raw := c.Body() // Ambil raw JSON body
	if len(raw) > 0 && raw[0] == '[' {
		// JSON berupa array
		if err := c.BodyParser(&multiple); err != nil {
			log.Printf("[ERROR] Failed to parse tooltip array: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON array"})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Received empty tooltip array")
			return c.Status(400).JSON(fiber.Map{"error": "Tooltip array is empty"})
		}

		// Validasi setiap item
		for i, tip := range multiple {
			if tip.Keyword == "" || tip.DescriptionShort == "" || tip.DescriptionLong == "" {
				log.Printf("[ERROR] Invalid tooltip at index %d: %+v\n", i, tip)
				return c.Status(400).JSON(fiber.Map{
					"error": "Each tooltip must have keyword, description_short, and description_long",
					"index": i,
					"data":  tip,
				})
			}
		}

		// Insert batch
		if err := tc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple tooltips: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create tooltips"})
		}

		log.Printf("[SUCCESS] Inserted %d tooltips", len(multiple))
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple tooltips created successfully",
			"data":    multiple,
		})
	}

	// Fallback: parse single object
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single tooltip: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format (expected object or array)",
		})
	}

	log.Printf("[DEBUG] Parsed single tooltip: %+v", single)

	if single.Keyword == "" || single.DescriptionShort == "" || single.DescriptionLong == "" {
		return c.Status(400).JSON(fiber.Map{"error": "keyword, description_short, and description_long are required"})
	}

	if err := tc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert tooltip: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tooltip"})
	}

	log.Printf("[SUCCESS] Tooltip created: ID=%d, Keyword=%s", single.ID, single.Keyword)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tooltip created successfully",
		"data":    single,
	})
}

// GetAllTooltips menangani permintaan untuk mendapatkan semua data tooltips
func (tc *TooltipsController) GetAllTooltips(c *fiber.Ctx) error {
	log.Println("Fetching all tooltips")

	var tooltips []model.Tooltip

	// Ambil semua data dari database
	if err := tc.DB.Find(&tooltips).Error; err != nil {
		log.Println("Error fetching tooltips:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tooltips"})
	}

	return c.JSON(tooltips)
}

func (tc *TooltipsController) UpdateTooltip(c *fiber.Ctx) error {
	id := c.Params("id")

	var existing model.Tooltip
	if err := tc.DB.First(&existing, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tooltip not found",
		})
	}

	var updated model.Tooltip
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Update fields
	existing.Keyword = updated.Keyword
	existing.DescriptionShort = updated.DescriptionShort
	existing.DescriptionLong = updated.DescriptionLong

	if err := tc.DB.Save(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tooltip",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Tooltip with ID %v updated successfully", existing.ID),
		"data":    existing,
	})
}

func (tc *TooltipsController) DeleteTooltip(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting tooltip with ID: %s\n", id)

	var tooltip model.Tooltip
	if err := tc.DB.First(&tooltip, id).Error; err != nil {
		log.Printf("[ERROR] Tooltip not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tooltip not found",
		})
	}

	if err := tc.DB.Delete(&tooltip).Error; err != nil {
		log.Printf("[ERROR] Failed to delete tooltip: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tooltip",
		})
	}

	log.Printf("[SUCCESS] Tooltip with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Tooltip with ID %s deleted successfully", id),
	})
}
