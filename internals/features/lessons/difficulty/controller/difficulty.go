package controller

import (
	"encoding/json"
	"fmt"
	"log"
	dto "quizku/internals/features/lessons/difficulty/dto"
	"quizku/internals/features/lessons/difficulty/model"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DifficultyController struct {
	DB *gorm.DB
}

func NewDifficultyController(db *gorm.DB) *DifficultyController {
	return &DifficultyController{DB: db}
}

// Get all difficulties
func (dc *DifficultyController) GetDifficulties(c *fiber.Ctx) error {
	var difficulties []model.DifficultyModel
	log.Println("[INFO] Received request to fetch all difficulties")

	if err := dc.DB.Find(&difficulties).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch difficulties: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var responses []dto.DifficultyResponse
	for _, d := range difficulties {
		responses = append(responses, dto.DifficultyResponse{
			DifficultyID:               d.DifficultyID,
			DifficultyName:             d.DifficultyName,
			DifficultyStatus:           d.DifficultyStatus,
			DifficultyDescriptionShort: d.DifficultyDescriptionShort,
			DifficultyDescriptionLong:  d.DifficultyDescriptionLong,
			DifficultyTotalCategories:  convertInt64ArrayToInt(d.DifficultyTotalCategories),
			DifficultyImageURL:         d.DifficultyImageURL,
		})
	}

	log.Printf("[SUCCESS] Retrieved %d difficulties\n", len(responses))
	return c.JSON(fiber.Map{
		"message": "All difficulties fetched successfully",
		"total":   len(responses),
		"data":    responses,
	})
}

// Get difficulty by ID
func (dc *DifficultyController) GetDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching difficulty with ID: %s\n", id)

	var difficulty model.DifficultyModel
	if err := dc.DB.First(&difficulty, id).Error; err != nil {
		log.Printf("[ERROR] Difficulty with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Difficulty not found"})
	}

	log.Printf("[SUCCESS] Retrieved difficulty: ID=%d, Name=%s\n", difficulty.DifficultyID, difficulty.DifficultyName)
	return c.JSON(fiber.Map{
		"message": "Difficulty fetched successfully",
		"data":    difficulty,
	})
}

// CreateDifficulty menangani input satu atau banyak difficulty
// CreateDifficulty menangani input satu atau banyak difficulty
func (dc *DifficultyController) CreateDifficulty(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create difficulty")

	var single model.DifficultyModel
	var multiple []model.DifficultyModel

	// 🧠 Coba parsing array terlebih dahulu
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d difficulties as array\n", len(multiple))

		for i, d := range multiple {
			if d.DifficultyName == "" {
				return c.Status(400).JSON(fiber.Map{
					"error": "difficulty_name is required in each difficulty",
					"index": i,
				})
			}
		}

		if err := dc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple difficulties: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create difficulties"})
		}

		log.Printf("[SUCCESS] %d difficulties created successfully\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple difficulties created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Jika bukan array, parse sebagai satu objek
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single difficulty input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("[DEBUG] Parsed single difficulty: %+v\n", single)

	if single.DifficultyName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "difficulty_name is required"})
	}

	if err := dc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert single difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create difficulty"})
	}

	log.Printf("[SUCCESS] Difficulty created: ID=%d, Name=%s\n", single.DifficultyID, single.DifficultyName)
	return c.Status(201).JSON(fiber.Map{
		"message": "Difficulty created successfully",
		"data":    single,
	})
}

// Update difficulty
func (dc *DifficultyController) UpdateDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating difficulty with ID: %s\n", id)

	var difficulty model.DifficultyModel
	if err := dc.DB.First(&difficulty, id).Error; err != nil {
		log.Printf("[ERROR] Difficulty with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Difficulty not found"})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Invalid JSON input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 🔄 Khusus untuk JSONB: difficulty_update_news
	if un, ok := input["difficulty_update_news"]; ok {
		jsonData, err := json.Marshal(un)
		if err == nil {
			input["difficulty_update_news"] = datatypes.JSON(jsonData)
		}
	}

	if err := dc.DB.Model(&difficulty).Updates(input).Error; err != nil {
		log.Printf("[ERROR] Failed to update difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUCCESS] Difficulty with ID %s updated successfully\n", id)
	return c.JSON(fiber.Map{
		"message": "Difficulty updated successfully",
		"data":    difficulty,
	})
}

// Delete difficulty
// 🟢 DELETE DIFFICULTY: Hapus difficulty berdasarkan ID (soft delete)
func (dc *DifficultyController) DeleteDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting difficulty with ID:", id)

	var difficulty model.DifficultyModel

	// 🔍 Cari difficulty berdasarkan ID
	if err := dc.DB.Where("difficulty_id = ?", id).First(&difficulty).Error; err != nil {
		log.Println("[ERROR] Difficulty tidak ditemukan:", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Difficulty tidak ditemukan",
		})
	}

	// 🗑️ Soft delete: menggunakan kolom deleted_at
	if err := dc.DB.Delete(&difficulty).Error; err != nil {
		log.Println("[ERROR] Gagal menghapus difficulty:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghapus difficulty",
		})
	}

	log.Printf("[SUCCESS] Difficulty dengan ID %s berhasil dihapus", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Difficulty dengan ID %s berhasil dihapus", id),
	})
}

func convertInt64ArrayToInt(arr pq.Int64Array) []int {
	var result []int
	for _, v := range arr {
		result = append(result, int(v))
	}
	return result
}
