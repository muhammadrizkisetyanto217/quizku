package controller

import (
	"fmt"
	"log"

	categoryModel "quizku/internals/features/lessons/categories/model"

	"quizku/internals/features/lessons/subcategories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubcategoryController struct {
	DB *gorm.DB
}

func NewSubcategoryController(db *gorm.DB) *SubcategoryController {
	return &SubcategoryController{DB: db}
}

// 🟢 GET ALL SUBCATEGORIES: Ambil seluruh data subkategori
func (sc *SubcategoryController) GetSubcategories(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all subcategories")
	var subcategories []model.SubcategoryModel

	// 🔍 Query semua subkategori
	if err := sc.DB.Find(&subcategories).Error; err != nil {
		log.Println("[ERROR] Failed to fetch subcategories:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}

	// ✅ Kirim data subkategori
	log.Printf("[SUCCESS] Retrieved %d subcategories\n", len(subcategories))
	return c.JSON(fiber.Map{
		"message": "All subcategories fetched successfully",
		"total":   len(subcategories),
		"data":    subcategories,
	})
}

// 🟢 GET SUBCATEGORY BY ID: Ambil data subkategori berdasarkan ID
func (sc *SubcategoryController) GetSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Query berdasarkan ID
	if err := sc.DB.First(&subcategory, id).Error; err != nil {
		log.Println("[ERROR] Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}

	// ✅ Kirim data
	log.Printf("[SUCCESS] Subcategory retrieved: ID=%d, Name=%s\n", subcategory.ID, subcategory.Name)
	return c.JSON(fiber.Map{
		"message": "Subcategory fetched successfully",
		"data":    subcategory,
	})
}

// 🟢 GET SUBCATEGORIES BY CATEGORY ID: Ambil data subkategori berdasarkan categories_id
func (sc *SubcategoryController) GetSubcategoriesByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	log.Printf("[INFO] Fetching subcategories with category ID: %s\n", categoryID)

	var subcategories []model.SubcategoryModel

	// 🔍 Query subkategori berdasarkan categories_id
	if err := sc.DB.Where("categories_id = ?", categoryID).Find(&subcategories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch subcategories for category ID %s: %v\n", categoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}

	// ✅ Kirim hasil
	log.Printf("[SUCCESS] Retrieved %d subcategories for category ID %s\n", len(subcategories), categoryID)
	return c.JSON(fiber.Map{
		"message": "Subcategories fetched successfully by category",
		"total":   len(subcategories),
		"data":    subcategories,
	})
}

// 🟢 CREATE SUBCATEGORY: Tambah satu atau banyak subkategori dengan validasi
func (sc *SubcategoryController) CreateSubcategory(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create subcategory")

	var single model.SubcategoryModel
	var multiple []model.SubcategoryModel

	// 🧠 Coba parse sebagai array
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d subcategories as array\n", len(multiple))

		// ✅ Validasi tiap item
		for i, item := range multiple {
			if item.Name == "" || item.CategoriesID == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "All fields are required in array (name, categories_id)",
					"index": i,
				})
			}

			var count int64
			if err := sc.DB.Table("categories").Where("id = ?", item.CategoriesID).Count(&count).Error; err != nil || count == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid categories_id in array",
					"index": i,
				})
			}
		}

		// 💾 Simpan semua jika valid
		if err := sc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple subcategories: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create subcategories"})
		}

		log.Printf("[SUCCESS] %d subcategories created successfully\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple subcategories created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Parse sebagai objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single subcategory input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	log.Printf("[DEBUG] Parsed single subcategory: %+v\n", single)

	// ✅ Validasi
	if single.Name == "" || single.CategoriesID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required (name, categories_id)"})
	}

	var count int64
	if err := sc.DB.Table("categories").Where("id = ?", single.CategoriesID).Count(&count).Error; err != nil || count == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid categories_id. Referenced category does not exist"})
	}

	// 💾 Simpan
	if err := sc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Error creating subcategory: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create subcategory"})
	}

	log.Printf("[SUCCESS] Subcategory created: ID=%d, Name=%s\n", single.ID, single.Name)
	return c.Status(201).JSON(fiber.Map{
		"message": "Subcategory created successfully",
		"data":    single,
	})
}

// 🟢 UPDATE SUBCATEGORY: Perbarui subkategori berdasarkan ID
func (sc *SubcategoryController) UpdateSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Cari berdasarkan ID
	if err := sc.DB.First(&subcategory, id).Error; err != nil {
		log.Println("[ERROR] Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}

	// 🔄 Update isi
	if err := c.BodyParser(&subcategory); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 💾 Simpan update
	if err := sc.DB.Save(&subcategory).Error; err != nil {
		log.Println("[ERROR] Error updating subcategory:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update subcategory"})
	}

	log.Printf("[SUCCESS] Subcategory updated: ID=%d, Name=%s\n", subcategory.ID, subcategory.Name)
	return c.JSON(fiber.Map{
		"message": "Subcategory updated successfully",
		"data":    subcategory,
	})
}

// 🟢 DELETE SUBCATEGORY: Hapus subkategori berdasarkan ID
func (sc *SubcategoryController) DeleteSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Cek apakah subkategori ada
	if err := sc.DB.First(&subcategory, id).Error; err != nil {
		log.Println("[ERROR] Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Subcategory not found",
		})
	}

	// 🗑️ Hapus data
	if err := sc.DB.Delete(&subcategory).Error; err != nil {
		log.Println("[ERROR] Error deleting subcategory:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete subcategory",
		})
	}

	// ✅ Kirim respon sukses
	log.Printf("[SUCCESS] Subcategory with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Subcategory with ID %s deleted successfully", id),
	})
}

// 🟢 GET CATEGORY WITH SUBCATEGORY AND THEMES: Ambil data lengkap kategori, subkategori, dan themes berdasarkan difficulty_id
func (sc *SubcategoryController) GetCategoryWithSubcategoryAndThemes(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Printf("[INFO] Fetching category, subcategory, and themes for difficulty ID: %s\n", difficultyID)

	// 🔍 Ambil semua kategori dengan relasi subkategori aktif & themes-nya
	var categories []categoryModel.CategoryModel
	if err := sc.DB.
		Where("difficulty_id = ?", difficultyID).
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("status = ?", "active").
				Preload("ThemesOrLevels")
		}).
		Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data kategori",
		})
	}

	// ✅ Kirim data lengkap
	log.Printf("[SUCCESS] Retrieved %d categories with subcategories and themes for difficulty ID %s\n", len(categories), difficultyID)
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data kategori lengkap",
		"data":    categories,
	})
}
