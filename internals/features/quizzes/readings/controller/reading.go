package controller

import (
	"fmt"
	readingModel "quizku/internals/features/quizzes/readings/model"
	tooltipModel "quizku/internals/features/utils/tooltips/model"

	"log"
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReadingController struct {
	DB *gorm.DB
}

func NewReadingController(db *gorm.DB) *ReadingController {
	return &ReadingController{DB: db}
}

// ✅ GET /api/readings
// Mengambil semua data bacaan yang tersedia di database.
func (rc *ReadingController) GetReadings(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all readings")

	var readings []readingModel.ReadingModel
	if err := rc.DB.Find(&readings).Error; err != nil {
		log.Println("[ERROR] Failed to fetch readings:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch readings",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d readings\n", len(readings))
	return c.JSON(readings)
}

// ✅ GET /api/readings/:id
// Mengambil satu data bacaan berdasarkan ID-nya.
func (rc *ReadingController) GetReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching reading with ID:", id)

	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, "reading_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Reading not found",
		})
	}

	return c.JSON(reading)
}

// ✅ GET /api/readings/unit/:unitId
// Mengambil semua data bacaan berdasarkan unit_id tertentu.
func (rc *ReadingController) GetReadingsByUnit(c *fiber.Ctx) error {
	unitID := c.Params("unitId")
	log.Printf("[INFO] Fetching readings for reading_unit_id: %s\n", unitID)

	var readings []readingModel.ReadingModel
	if err := rc.DB.Where("reading_unit_id = ?", unitID).Find(&readings).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch readings for reading_unit_id %s: %v\n", unitID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch readings",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d readings for reading_unit_id %s\n", len(readings), unitID)
	return c.JSON(readings)
}

// ✅ POST /api/readings
// Membuat satu data bacaan baru
func (rc *ReadingController) CreateReading(c *fiber.Ctx) error {
	log.Println("[INFO] Creating a new reading")

	var reading readingModel.ReadingModel
	if err := c.BodyParser(&reading); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := rc.DB.Create(&reading).Error; err != nil {
		log.Println("[ERROR] Failed to create reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create reading",
		})
	}

	log.Printf("[SUCCESS] Reading created: ID=%d\n", reading.ReadingID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reading created successfully",
		"data":    reading,
	})
}

// ✅ PUT /api/readings/:id
// Memperbarui data bacaan berdasarkan ID
func (rc *ReadingController) UpdateReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating reading with ID: %s\n", id)

	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, "reading_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Reading not found",
		})
	}

	if err := c.BodyParser(&reading); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := rc.DB.Save(&reading).Error; err != nil {
		log.Println("[ERROR] Failed to update reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update reading",
		})
	}

	log.Printf("[SUCCESS] Reading updated: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Reading updated successfully",
		"data":    reading,
	})
}

// ✅ DELETE /api/readings/:id
// Menghapus data bacaan berdasarkan ID
func (rc *ReadingController) DeleteReading(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting reading with ID: %s\n", id)

	if err := rc.DB.Delete(&readingModel.ReadingModel{}, "reading_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Failed to delete reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete reading",
		})
	}

	log.Printf("[SUCCESS] Reading with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Reading with ID %s deleted successfully", id),
	})
}

// GetReadingWithTooltips digunakan untuk mengambil satu data reading berdasarkan ID.
// Saat ini fungsinya hanya mengambil data dari database dan belum memodifikasi teks reading.
// Jika ingin menambahkan tooltip ke dalam teks, perlu menggabungkan fungsi ini dengan MarkKeywords()
// untuk menyisipkan ID tooltip ke dalam teks bacaan.

func (rc *ReadingController) GetReadingWithTooltips(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching reading with ID: %s\n", id)

	var reading readingModel.ReadingModel
	if err := rc.DB.First(&reading, "reading_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Reading not found",
		})
	}

	log.Printf("[SUCCESS] Retrieved reading with ID: %s\n", id)
	return c.JSON(fiber.Map{
		"reading": reading,
	})
}

// MarkKeywords menandai setiap kemunculan keyword dalam teks dengan format tambahan "=ID".
// Contoh: jika keyword adalah "shalat" dengan ID 3, maka dalam teks akan muncul sebagai "shalat=3".
// - Menggunakan regex untuk pencocokan kata secara case-insensitive.
// - Preserve casing asli dari teks (tidak diubah ke huruf kecil).
// - Aman dari konflik keyword karena memakai regexp.QuoteMeta untuk escape karakter khusus.
// Fungsi ini cocok dipakai untuk menandai teks yang akan diolah di frontend sebagai tooltip atau definisi.

func (rc *ReadingController) MarkKeywords(text string, tooltips []tooltipModel.Tooltip) string {
	log.Printf("[DEBUG] Original text: %s\n", text)

	for _, tooltip := range tooltips {
		keyword := tooltip.TooltipKeyword
		keywordID := strconv.Itoa(int(tooltip.TooltipID))

		// Regex case-insensitive tapi preserve original match
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return match + "=" + keywordID
		})

		log.Printf("[DEBUG] Replacing all '%s' with '%s' in text", keyword, keyword+"="+keywordID)
	}

	log.Printf("[DEBUG] Modified text: %s\n", text)
	return text
}