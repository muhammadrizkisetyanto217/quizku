package controller

import (
	"fmt"
	readingModel "quizku/internals/features/quizzes/readings/model"
	tooltipModel "quizku/internals/features/utils/tooltips/model"

	"log"
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ReadingController struct {
	DB *gorm.DB
}

func NewReadingController(db *gorm.DB) *ReadingController {
	return &ReadingController{DB: db}
}

// Get all readings
func (rc *ReadingController) GetReadings(c *fiber.Ctx) error {
	log.Println("Fetching all readings")
	var readings []readingModel.ReadingModel
	if err := rc.DB.Find(&readings).Error; err != nil {
		log.Println("Error fetching readings:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch readings"})
	}
	return c.JSON(readings)
}

// Get a single reading by ID
func (rc *ReadingController) GetReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Fetching reading with ID:", id)
	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, id).Error; err != nil {
		log.Println("Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reading not found"})
	}
	return c.JSON(reading)
}

// Get readings by unit ID
func (rc *ReadingController) GetReadingsByUnit(c *fiber.Ctx) error {
	unitID := c.Params("unitId")
	log.Printf("[INFO] Fetching readings for unit_id: %s\n", unitID)

	var readings []readingModel.ReadingModel
	if err := rc.DB.Where("unit_id = ?", unitID).Find(&readings).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch readings for unit_id %s: %v\n", unitID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch readings"})
	}
	log.Printf("[SUCCESS] Retrieved %d readings for unit_id %s\n", len(readings), unitID)
	return c.JSON(readings)
}

// Create a new reading
func (rc *ReadingController) CreateReading(c *fiber.Ctx) error {
	log.Println("Creating a new reading")
	var reading readingModel.ReadingModel
	if err := c.BodyParser(&reading); err != nil {
		log.Println("Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := rc.DB.Create(&reading).Error; err != nil {
		log.Println("Error creating reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create reading"})
	}
	return c.Status(fiber.StatusCreated).JSON(reading)
}

// Update a reading
func (rc *ReadingController) UpdateReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Updating reading with ID:", id)
	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, id).Error; err != nil {
		log.Println("Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reading not found"})
	}
	if err := c.BodyParser(&reading); err != nil {
		log.Println("Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := rc.DB.Save(&reading).Error; err != nil {
		log.Println("Error updating reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update reading"})
	}
	return c.JSON(reading)
}

// Delete a reading
func (rc *ReadingController) DeleteReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Deleting reading with ID:", id)

	if err := rc.DB.Delete(&readingModel.ReadingModel{}, id).Error; err != nil {
		log.Println("Error deleting reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete reading",
		})
	}

	log.Printf("[SUCCESS] Reading with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Reading with ID %s deleted successfully", id),
	})
}

// Get a single reading by ID with Tooltips
func (rc *ReadingController) GetReadingWithTooltips(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching reading with ID: %s\n", id)

	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reading not found"})
	}

	// Fetch Tooltips
	var tooltips []tooltipModel.Tooltip
	if len(reading.TooltipsID) > 0 {
		if err := rc.DB.Where("id = ANY(?)", pq.Array(reading.TooltipsID)).Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tooltips"})
		}
	}

	log.Printf("[SUCCESS] Retrieved reading with ID: %s\n", id)
	return c.JSON(fiber.Map{
		"reading":  reading,
		"tooltips": tooltips,
	})
}

// Get a onlyReading by ID with Tooltips
func (rc *ReadingController) GetOnlyReadingTooltips(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching reading with ID: %s\n", id)

	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reading not found"})
	}

	// Fetch Tooltips
	var tooltips []tooltipModel.Tooltip
	if len(reading.TooltipsID) > 0 {
		if err := rc.DB.Where("id = ANY(?)", pq.Array(reading.TooltipsID)).Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tooltips"})
		}
	}

	log.Printf("[SUCCESS] Retrieved reading with ID: %s\n", id)
	return c.JSON(fiber.Map{
		// "reading":  reading,
		"tooltips": tooltips,
	})
}

func (rc *ReadingController) MarkKeywords(text string, tooltips []tooltipModel.Tooltip) string {
	log.Printf("[DEBUG] Original text: %s\n", text)

	for _, tooltip := range tooltips {
		keyword := tooltip.Keyword
		keywordID := strconv.Itoa(int(tooltip.ID))

		// Regex case-insensitive tapi preserve original match
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return match + "=" + keywordID // Tetap gunakan `match` agar case aslinya dipertahankan
		})

		log.Printf("[DEBUG] Replacing all '%s' with '%s' in text", keyword, keyword+"="+keywordID)
	}

	log.Printf("[DEBUG] Modified text: %s\n", text)
	return text
}

// **📌 Get Reading by ID dengan Tooltips yang Ditandai dan Update ke Database**
func (rc *ReadingController) ConvertReadingWithTooltipsId(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching reading with ID: %s\n", id)

	// **📌 Ambil Data Reading**
	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reading not found"})
	}

	// **📌 Ambil Tooltips yang Sesuai**
	var tooltips []tooltipModel.Tooltip
	if len(reading.TooltipsID) > 0 {
		query := rc.DB.Where("id = ANY(?)", pq.Array(reading.TooltipsID))
		if err := query.Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tooltips"})
		}
	}

	// **📌 Tandai Keyword dalam Title & Description (HINDARI DUPLIKASI)**
	markedTitle := rc.MarkKeywords(reading.Title, tooltips)
	markedDescription := rc.MarkKeywords(reading.DescriptionLong, tooltips)

	// **📌 Hanya update jika ada perubahan**
	if markedTitle != reading.Title || markedDescription != reading.DescriptionLong {
		if err := rc.DB.Model(&reading).Updates(map[string]interface{}{
			"title":            markedTitle,
			"description_long": markedDescription,
		}).Error; err != nil {
			log.Println("[ERROR] Failed to update reading:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update reading"})
		}
		log.Printf("[SUCCESS] Updated reading with ID: %s\n", id)
	} else {
		log.Printf("[INFO] No changes detected, skipping update for ID: %s\n", id)
	}

	// **📌 Kembalikan Response**
	return c.JSON(fiber.Map{
		"reading": fiber.Map{
			"id":               reading.ID,
			"title":            markedTitle,
			"status":           reading.Status,
			"description_long": markedDescription,
			"tooltips_id":      reading.TooltipsID,
			"created_at":       reading.CreatedAt,
			"updated_at":       reading.UpdatedAt,
			"deleted_at":       reading.DeletedAt,
			"unit_id":          reading.UnitID,
			"created_by":       reading.CreatedBy,
		},
		"tooltips": tooltips,
	})
}
